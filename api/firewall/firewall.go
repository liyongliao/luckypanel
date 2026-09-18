package firewall

import (
	"net/http"
	"strconv"

	"github.com/0xJacky/Nginx-UI/internal/firewall"
	"github.com/0xJacky/Nginx-UI/model"
	"github.com/gin-gonic/gin"
)

// GetStatus returns current firewall daemon status.
func GetStatus(c *gin.Context) {
	status := firewall.GetStatus()
	c.JSON(http.StatusOK, status)
}

// ToggleStatus enables/disables the firewall.
func ToggleStatus(c *gin.Context) {
	var body struct {
		Enable bool `json:"enable"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := firewall.SetStatus(body.Enable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetRules returns all saved firewall rules.
func GetRules(c *gin.Context) {
	db := model.UseDB()
	var rules []model.FirewallRule
	if err := db.Order("id desc").Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

// OpenPort opens a new port in the firewall.
func OpenPort(c *gin.Context) {
	var body struct {
		Port        string `json:"port" binding:"required"`
		Protocol    string `json:"protocol"`
		Source      string `json:"source"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := firewall.OpenPort(body.Port, body.Protocol, body.Source, body.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// ClosePort removes a port rule by ID.
func ClosePort(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	if err := firewall.ClosePort(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
