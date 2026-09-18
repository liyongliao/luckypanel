package firewall

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.RouterGroup) {
	fw := r.Group("firewall")
	{
		fw.GET("status", GetStatus)
		fw.POST("toggle", ToggleStatus)
		fw.GET("rules", GetRules)
		fw.POST("rules", OpenPort)
		fw.DELETE("rules/:id", ClosePort)
	}
}
