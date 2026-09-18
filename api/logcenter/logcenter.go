package logcenter

import (
	"net/http"
	"strconv"

	internalLog "github.com/0xJacky/Nginx-UI/internal/logcenter"
	"github.com/gin-gonic/gin"
)

// GetSources returns all queryable log sources.
func GetSources(c *gin.Context) {
	sources := internalLog.GetSources()
	c.JSON(http.StatusOK, sources)
}

// QueryLogs searches logs from a source.
func QueryLogs(c *gin.Context) {
	source := c.DefaultQuery("source", "nginx_access")
	keyword := c.DefaultQuery("keyword", "")
	limitStr := c.DefaultQuery("limit", "200")
	limit, _ := strconv.Atoi(limitStr)

	lines, err := internalLog.QueryLogs(source, keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"source":  source,
		"keyword": keyword,
		"count":   len(lines),
		"lines":   lines,
	})
}
