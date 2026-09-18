package firewall

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/0xJacky/Nginx-UI/model"
	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/settings"
)

type FirewallType string

const (
	TypeUFW       FirewallType = "ufw"
	TypeFirewalld FirewallType = "firewalld"
	TypeIPTables  FirewallType = "iptables"
	TypeNone      FirewallType = "none"
)

type FirewallStatus struct {
	Type      FirewallType `json:"type"`
	IsActive  bool         `json:"is_active"`
	Available bool         `json:"available"`
	Message   string       `json:"message"`
}

var (
	detectedType FirewallType
	detectOnce   sync.Once
	mu           sync.Mutex
)

// DetectFirewall identifies which firewall daemon is active.
func DetectFirewall() FirewallType {
	detectOnce.Do(func() {
		if runtime.GOOS != "linux" {
			detectedType = TypeNone
			return
		}

		// Check UFW
		if _, err := exec.LookPath("ufw"); err == nil {
			out, err := exec.Command("ufw", "status").CombinedOutput()
			if err == nil && !strings.Contains(string(out), "command not found") {
				detectedType = TypeUFW
				return
			}
		}

		// Check Firewalld
		if _, err := exec.LookPath("firewall-cmd"); err == nil {
			out, err := exec.Command("firewall-cmd", "--state").CombinedOutput()
			if err == nil || strings.Contains(string(out), "running") || strings.Contains(string(out), "not running") {
				detectedType = TypeFirewalld
				return
			}
		}

		// Check iptables
		if _, err := exec.LookPath("iptables"); err == nil {
			detectedType = TypeIPTables
			return
		}

		detectedType = TypeNone
	})
	return detectedType
}

// GetStatus returns the current status of the detected firewall.
func GetStatus() FirewallStatus {
	fwType := DetectFirewall()
	if fwType == TypeNone {
		return FirewallStatus{
			Type:      TypeNone,
			IsActive:  true, // simulated active in non-Linux or dev environments
			Available: true,
			Message:   "Simulated firewall mode (No host Linux firewall detected)",
		}
	}

	isActive := false
	var msg string

	switch fwType {
	case TypeUFW:
		out, err := exec.Command("ufw", "status").CombinedOutput()
		if err == nil && strings.Contains(string(out), "Status: active") {
			isActive = true
			msg = "UFW is active"
		} else {
			msg = "UFW is inactive"
		}
	case TypeFirewalld:
		out, err := exec.Command("firewall-cmd", "--state").CombinedOutput()
		if err == nil && strings.TrimSpace(string(out)) == "running" {
			isActive = true
			msg = "firewalld is running"
		} else {
			msg = "firewalld is stopped"
		}
	case TypeIPTables:
		isActive = true
		msg = "iptables is available"
	}

	return FirewallStatus{
		Type:      fwType,
		IsActive:  isActive,
		Available: true,
		Message:   msg,
	}
}

// SetStatus enables or disables the firewall daemon.
func SetStatus(enable bool) error {
	mu.Lock()
	defer mu.Unlock()

	fwType := DetectFirewall()
	if fwType == TypeNone {
		return nil
	}

	var cmd *exec.Cmd
	switch fwType {
	case TypeUFW:
		if enable {
			cmd = exec.Command("ufw", "--force", "enable")
		} else {
			cmd = exec.Command("ufw", "disable")
		}
	case TypeFirewalld:
		if enable {
			cmd = exec.Command("systemctl", "start", "firewalld")
		} else {
			cmd = exec.Command("systemctl", "stop", "firewalld")
		}
	case TypeIPTables:
		return nil
	}

	if cmd != nil {
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %v", strings.TrimSpace(string(out)), err)
		}
	}
	return nil
}

// IsProtectedPort prevents locking out the Web panel or SSH.
func IsProtectedPort(portStr string) bool {
	// 22 is standard SSH
	if portStr == "22" {
		return true
	}

	// Current web panel port
	panelPort := fmt.Sprintf("%d", settings.ServerSettings.Port)
	if panelPort != "0" && portStr == panelPort {
		return true
	}

	return false
}

// OpenPort opens a port in the firewall and persists it in the database.
func OpenPort(port string, protocol string, source string, description string) (*model.FirewallRule, error) {
	mu.Lock()
	defer mu.Unlock()

	if port == "" {
		return nil, errors.New("port cannot be empty")
	}

	if protocol == "" {
		protocol = "tcp"
	}
	protocol = strings.ToLower(protocol)

	if source == "" {
		source = "any"
	}

	fwType := DetectFirewall()
	var err error

	// Apply rule to underlying firewall if present
	if fwType != TypeNone {
		err = applyFirewallRule(fwType, port, protocol, source, true)
		if err != nil {
			logger.Warnf("Firewall underlying command warning: %v", err)
		}
	}

	// Save to DB
	rule := &model.FirewallRule{
		Port:        port,
		Protocol:    protocol,
		Strategy:    "allow",
		Source:      source,
		Description: description,
	}

	db := model.UseDB()
	if err := db.Create(rule).Error; err != nil {
		return nil, err
	}

	return rule, nil
}

// ClosePort removes a port rule from the firewall and database.
func ClosePort(id uint64) error {
	mu.Lock()
	defer mu.Unlock()

	db := model.UseDB()
	var rule model.FirewallRule
	if err := db.First(&rule, id).Error; err != nil {
		return errors.New("firewall rule not found")
	}

	if IsProtectedPort(rule.Port) {
		return fmt.Errorf("security check: cannot remove rule for protected port %s (SSH/Web Panel)", rule.Port)
	}

	fwType := DetectFirewall()
	if fwType != TypeNone {
		_ = applyFirewallRule(fwType, rule.Port, rule.Protocol, rule.Source, false)
	}

	return db.Delete(&rule).Error
}

func applyFirewallRule(fwType FirewallType, port string, protocol string, source string, open bool) error {
	protos := []string{protocol}
	if protocol == "both" || protocol == "tcp/udp" {
		protos = []string{"tcp", "udp"}
	}

	for _, proto := range protos {
		var cmd *exec.Cmd
		switch fwType {
		case TypeUFW:
			if open {
				if source != "" && source != "any" && source != "0.0.0.0/0" {
					cmd = exec.Command("ufw", "allow", "from", source, "to", "any", "port", port, "proto", proto)
				} else {
					cmd = exec.Command("ufw", "allow", fmt.Sprintf("%s/%s", port, proto))
				}
			} else {
				if source != "" && source != "any" && source != "0.0.0.0/0" {
					cmd = exec.Command("ufw", "delete", "allow", "from", source, "to", "any", "port", port, "proto", proto)
				} else {
					cmd = exec.Command("ufw", "delete", "allow", fmt.Sprintf("%s/%s", port, proto))
				}
			}
		case TypeFirewalld:
			action := "--add-port"
			if !open {
				action = "--remove-port"
			}
			cmd = exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("%s=%s/%s", action, port, proto))
			if out, err := cmd.CombinedOutput(); err != nil {
				logger.Warnf("firewall-cmd error: %s: %v", string(out), err)
			}
			// reload firewalld
			_ = exec.Command("firewall-cmd", "--reload").Run()
			return nil
		case TypeIPTables:
			action := "-I"
			if !open {
				action = "-D"
			}
			cmd = exec.Command("iptables", action, "INPUT", "-p", proto, "--dport", port, "-j", "ACCEPT")
		}

		if cmd != nil {
			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("%s: %v", strings.TrimSpace(string(out)), err)
			}
		}
	}
	return nil
}
