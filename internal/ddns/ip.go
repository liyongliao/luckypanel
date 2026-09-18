package ddns

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

var defaultIPv4URLs = []string{
	"https://api.ipify.org",
	"https://ipv4.icanhazip.com",
	"https://ip.sb",
}

var defaultIPv6URLs = []string{
	"https://api64.ipify.org",
	"https://ipv6.icanhazip.com",
}

// GetIPByMethod retrieves IPv4 or IPv6 based on the chosen method.
func GetIPByMethod(ipType, method, targetNic, customURL string) (string, error) {
	switch method {
	case "nic":
		return GetIPFromInterface(targetNic, ipType)
	case "url":
		return GetIPFromURL(customURL, ipType)
	default:
		return GetIPFromURL(customURL, ipType)
	}
}

// GetIPFromURL fetches public IP from online probes.
func GetIPFromURL(customURL, ipType string) (string, error) {
	urls := defaultIPv4URLs
	if ipType == "ipv6" {
		urls = defaultIPv6URLs
	}
	if customURL != "" {
		urls = []string{customURL}
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, u := range urls {
		req, err := http.NewRequestWithContext(context.Background(), "GET", u, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		ipStr := strings.TrimSpace(string(body))
		parsedIP := net.ParseIP(ipStr)
		if parsedIP != nil {
			if ipType == "ipv4" && parsedIP.To4() != nil {
				return ipStr, nil
			}
			if ipType == "ipv6" && parsedIP.To4() == nil {
				return ipStr, nil
			}
		}
	}

	return "", errors.New("failed to acquire IP from probe URLs")
}

// GetIPFromInterface extracts IPv4 or IPv6 from a specific local network card.
func GetIPFromInterface(nicName, ipType string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if nicName != "" && iface.Name != nicName {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			if ipType == "ipv4" {
				if ip4 := ip.To4(); ip4 != nil {
					return ip4.String(), nil
				}
			} else if ipType == "ipv6" {
				if ip.To4() == nil && !ip.IsLinkLocalUnicast() {
					return ip.String(), nil
				}
			}
		}
	}

	return "", errors.New("no matching IP found on interface: " + nicName)
}

type InterfaceInfo struct {
	Name  string   `json:"name"`
	IPv4s []string `json:"ipv4s"`
	IPv6s []string `json:"ipv6s"`
}

// ListInterfaces returns all network cards and their current IPs for frontend dropdown selection.
func ListInterfaces() ([]InterfaceInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []InterfaceInfo
	for _, iface := range ifaces {
		info := InterfaceInfo{
			Name:  iface.Name,
			IPv4s: []string{},
			IPv6s: []string{},
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				info.IPv4s = append(info.IPv4s, ip4.String())
			} else if !ip.IsLinkLocalUnicast() {
				info.IPv6s = append(info.IPv6s, ip.String())
			}
		}
		result = append(result, info)
	}

	return result, nil
}
