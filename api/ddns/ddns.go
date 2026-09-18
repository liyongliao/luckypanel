package ddns

import (
	"net/http"
	"strconv"

	internalDDNS "github.com/0xJacky/Nginx-UI/internal/ddns"
	"github.com/0xJacky/Nginx-UI/model"
	"github.com/gin-gonic/gin"
)

// GetTasks returns all DDNS tasks.
func GetTasks(c *gin.Context) {
	db := model.UseDB()
	var tasks []model.DDNSTask
	if err := db.Order("id desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// CreateTask adds a new DDNS task.
func CreateTask(c *gin.Context) {
	var task model.DDNSTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := model.UseDB()
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Trigger immediate initial run in background
	if task.Enabled {
		go internalDDNS.RunTask(task.ID)
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask modifies an existing DDNS task.
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	db := model.UseDB()
	var existing model.DDNSTask
	if err := db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	var input model.DDNSTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.Name = input.Name
	existing.Provider = input.Provider
	existing.Domains = input.Domains
	existing.IPType = input.IPType
	existing.IPMethod = input.IPMethod
	existing.IPInterface = input.IPInterface
	existing.IPUrl = input.IPUrl
	if input.AccessKeyID != "" {
		existing.AccessKeyID = input.AccessKeyID
	}
	if input.AccessKeySecret != "" {
		existing.AccessKeySecret = input.AccessKeySecret
	}
	existing.WebhookURL = input.WebhookURL
	existing.IntervalSeconds = input.IntervalSeconds
	existing.Enabled = input.Enabled

	if err := db.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existing.Enabled {
		go internalDDNS.RunTask(existing.ID)
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteTask removes a DDNS task.
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	db := model.UseDB()
	if err := db.Delete(&model.DDNSTask{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// RunTaskNow immediately triggers the DDNS task.
func RunTaskNow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	if err := internalDDNS.RunTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db := model.UseDB()
	var updated model.DDNSTask
	_ = db.First(&updated, id)
	c.JSON(http.StatusOK, updated)
}

// GetInterfaces returns available host network cards.
func GetInterfaces(c *gin.Context) {
	ifaces, err := internalDDNS.ListInterfaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ifaces)
}
