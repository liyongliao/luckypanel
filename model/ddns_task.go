package model

import "time"

// DDNSTask stores dynamic DNS synchronization job configuration.
type DDNSTask struct {
	ID              uint64     `gorm:"primary_key" json:"id"`
	Name            string     `gorm:"size:128;not null" json:"name"`
	Provider        string     `gorm:"size:64;not null" json:"provider"` // "aliyun", "tencent", "cloudflare", "huawei", "callback"
	Domains         string     `gorm:"type:text;not null" json:"domains"` // comma-separated domains, e.g. "sub.domain.com,*.domain.com"
	IPType          string     `gorm:"size:16;default:'ipv4'" json:"ip_type"` // "ipv4", "ipv6", "both"
	IPMethod        string     `gorm:"size:32;default:'url'" json:"ip_method"` // "url", "nic", "stun"
	IPInterface     string     `gorm:"size:64" json:"ip_interface"`
	IPUrl           string     `gorm:"size:255" json:"ip_url"`
	AccessKeyID     string     `gorm:"size:255" json:"access_key_id"`
	AccessKeySecret string     `gorm:"size:255" json:"access_key_secret"`
	WebhookURL      string     `gorm:"size:512" json:"webhook_url"`
	IntervalSeconds int        `gorm:"default:300" json:"interval_seconds"`
	Enabled         bool       `gorm:"default:true" json:"enabled"`
	LastIPv4        string     `gorm:"size:64" json:"last_ipv4"`
	LastIPv6        string     `gorm:"size:128" json:"last_ipv6"`
	LastRunAt       *time.Time `json:"last_run_at"`
	LastStatus      string     `gorm:"size:32" json:"last_status"` // "success", "failed"
	LastError       string     `gorm:"type:text" json:"last_error"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (DDNSTask) TableName() string {
	return "ddns_tasks"
}
