package model

import "time"

// FirewallRule represents an allowed/denied port rule in the firewall.
type FirewallRule struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	Port        string    `gorm:"size:64;not null" json:"port"`             // e.g. "80", "8000-8080"
	Protocol    string    `gorm:"size:16;not null;default:'tcp'" json:"protocol"` // "tcp", "udp", "both"
	Strategy    string    `gorm:"size:16;not null;default:'allow'" json:"strategy"` // "allow", "deny"
	Source      string    `gorm:"size:128;default:'any'" json:"source"`       // "any" or specific IP/CIDR
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (FirewallRule) TableName() string {
	return "firewall_rules"
}
