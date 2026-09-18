package forward

import (
	"net/http"
	"strconv"

	"github.com/0xJacky/Nginx-UI/internal/forward"
	"github.com/0xJacky/Nginx-UI/model"
	"github.com/gin-gonic/gin"
)

type RuleWithStats struct {
	model.PortForwardRule
	ActiveConns int64 `json:"active_conns"`
	IsRunning   bool  `json:"is_running"`
}

// GetRules lists all forwarding rules with realtime stats.
func GetRules(c *gin.Context) {
	db := model.UseDB()
	var rules []model.PortForwardRule
	if err := db.Order("id desc").Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]RuleWithStats, len(rules))
	for i, r := range rules {
		activeConns, rx, tx, running := forward.GlobalManager.GetInstanceStats(r.ID)
		if running {
			r.RxBytes = rx
			r.TxBytes = tx
		}
		result[i] = RuleWithStats{
			PortForwardRule: r,
			ActiveConns:     activeConns,
			IsRunning:       running,
		}
	}

	c.JSON(http.StatusOK, result)
}

// CreateRule creates and activates a new port forwarding rule.
func CreateRule(c *gin.Context) {
	var rule model.PortForwardRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if rule.ListenPort <= 0 || rule.TargetPort <= 0 || rule.TargetIP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ports or target host"})
		return
	}

	db := model.UseDB()
	if err := db.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rule.Enabled {
		if err := forward.GlobalManager.StartRule(&rule); err != nil {
			c.JSON(http.StatusOK, gin.H{"data": rule, "warning": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, rule)
}

// UpdateRule updates rule configuration.
func UpdateRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	db := model.UseDB()
	var existing model.PortForwardRule
	if err := db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	var input model.PortForwardRule
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.Name = input.Name
	existing.Protocol = input.Protocol
	existing.ListenIP = input.ListenIP
	existing.ListenPort = input.ListenPort
	existing.TargetIP = input.TargetIP
	existing.TargetPort = input.TargetPort
	existing.Enabled = input.Enabled
	existing.EnableUPnP = input.EnableUPnP
	existing.AutoOpenFirewall = input.AutoOpenFirewall
	existing.Description = input.Description

	if err := db.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existing.Enabled {
		_ = forward.GlobalManager.StartRule(&existing)
	} else {
		forward.GlobalManager.StopRule(existing.ID)
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteRule removes a rule and shuts down its listener.
func DeleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	forward.GlobalManager.StopRule(id)

	db := model.UseDB()
	if err := db.Delete(&model.PortForwardRule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ToggleRule enables or disables a rule.
func ToggleRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	db := model.UseDB()
	var rule model.PortForwardRule
	if err := db.First(&rule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	rule.Enabled = !rule.Enabled
	_ = db.Save(&rule)

	if rule.Enabled {
		if err := forward.GlobalManager.StartRule(&rule); err != nil {
			c.JSON(http.StatusOK, gin.H{"rule": rule, "warning": err.Error()})
			return
		}
	} else {
		forward.GlobalManager.StopRule(rule.ID)
	}

	c.JSON(http.StatusOK, rule)
}
