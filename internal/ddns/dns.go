package ddns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0xJacky/Nginx-UI/model"
	"github.com/uozi-tech/cosy/logger"
)

// UpdateRecord updates the DNS records for the specified task with the new IP.
func UpdateRecord(task *model.DDNSTask, ip string) error {
	domains := strings.Split(task.Domains, ",")
	for _, d := range domains {
		domain := strings.TrimSpace(d)
		if domain == "" {
			continue
		}

		recordType := "A"
		if strings.Contains(ip, ":") {
			recordType = "AAAA"
		}

		var err error
		switch strings.ToLower(task.Provider) {
		case "cloudflare":
			err = updateCloudflare(task.AccessKeyID, domain, recordType, ip)
		case "callback", "webhook":
			err = triggerWebhook(task.WebhookURL, domain, recordType, ip)
		default:
			// Fallback / simulation / custom webhook
			if task.WebhookURL != "" {
				err = triggerWebhook(task.WebhookURL, domain, recordType, ip)
			} else {
				logger.Infof("[DDNS] Simulated update for %s -> %s (%s) on provider %s", domain, ip, recordType, task.Provider)
			}
		}

		if err != nil {
			return fmt.Errorf("domain %s update failed: %w", domain, err)
		}
	}
	return nil
}

func updateCloudflare(apiToken, domain, recordType, ip string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	// Extract zone domain (e.g. sub.example.com -> example.com)
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid domain: %s", domain)
	}
	zoneName := strings.Join(parts[len(parts)-2:], ".")

	// 1. Get Zone ID
	zoneReq, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.cloudflare.com/client/v4/zones?name="+zoneName, nil)
	zoneReq.Header.Set("Authorization", "Bearer "+apiToken)
	zoneReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(zoneReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var zoneResp struct {
		Success bool `json:"success"`
		Result  []struct {
			ID string `json:"id"`
		} `json:"result"`
		Errors []any `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&zoneResp); err != nil || !zoneResp.Success || len(zoneResp.Result) == 0 {
		return fmt.Errorf("failed to fetch Cloudflare zone for %s", zoneName)
	}
	zoneID := zoneResp.Result[0].ID

	// 2. Search existing record
	recordReq, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=%s&name=%s", zoneID, recordType, domain), nil)
	recordReq.Header.Set("Authorization", "Bearer "+apiToken)
	recordReq.Header.Set("Content-Type", "application/json")

	recResp, err := client.Do(recordReq)
	if err != nil {
		return err
	}
	defer recResp.Body.Close()

	var dnsListResp struct {
		Success bool `json:"success"`
		Result  []struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		} `json:"result"`
	}
	_ = json.NewDecoder(recResp.Body).Decode(&dnsListResp)

	payload := map[string]interface{}{
		"type":    recordType,
		"name":    domain,
		"content": ip,
		"ttl":     1, // Auto
		"proxied": false,
	}
	jsonBytes, _ := json.Marshal(payload)

	var updateReq *http.Request
	if len(dnsListResp.Result) > 0 {
		recordID := dnsListResp.Result[0].ID
		updateReq, _ = http.NewRequestWithContext(context.Background(), "PUT", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID), bytes.NewBuffer(jsonBytes))
	} else {
		updateReq, _ = http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID), bytes.NewBuffer(jsonBytes))
	}

	updateReq.Header.Set("Authorization", "Bearer "+apiToken)
	updateReq.Header.Set("Content-Type", "application/json")

	upResp, err := client.Do(updateReq)
	if err != nil {
		return err
	}
	defer upResp.Body.Close()

	var finalResp struct {
		Success bool `json:"success"`
	}
	_ = json.NewDecoder(upResp.Body).Decode(&finalResp)
	if !finalResp.Success {
		return fmt.Errorf("Cloudflare record update rejected")
	}

	return nil
}

func triggerWebhook(urlTemplate, domain, recordType, ip string) error {
	if urlTemplate == "" {
		return nil
	}

	url := strings.ReplaceAll(urlTemplate, "#{ip}", ip)
	url = strings.ReplaceAll(url, "#{domain}", domain)
	url = strings.ReplaceAll(url, "#{type}", recordType)

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}
