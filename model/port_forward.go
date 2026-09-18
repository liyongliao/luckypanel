package model

import "time"

// PortForwardRule stores TCP/UDP stream forwarding configuration.
type PortForwardRule struct {
	ID               uint64    `gorm:"primary_key" json:"id"`
	Name             string    `gorm:"size:128;not null" json:"name"`
	Protocol         string    `gorm:"size:16;not null;default:'tcp'" json:"protocol"` // "tcp", "udp", "both"
	ListenIP         string    `gorm:"size:64;default:'0.0.0.0'" json:"listen_ip"`
	ListenPort       int       `gorm:"not null" json:"listen_port"`
	TargetIP         string    `gorm:"size:128;not null" json:"target_ip"`
	TargetPort       int       `gorm:"not null" json:"target_port"`
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	EnableUPnP       bool      `gorm:"default:false" json:"enable_upnp"`
	AutoOpenFirewall bool      `gorm:"default:true" json:"auto_open_firewall"`
	RxBytes          uint64    `gorm:"default:0" json:"rx_bytes"`
	TxBytes          uint64    `gorm:"default:0" json:"tx_bytes"`
	Description      string    `gorm:"size:255" json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (PortForwardRule) TableName() string {
	return "port_forward_rules"
}
