package forward

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.RouterGroup) {
	fw := r.Group("forward")
	{
		fw.GET("rules", GetRules)
		fw.POST("rules", CreateRule)
		fw.PUT("rules/:id", UpdateRule)
		fw.DELETE("rules/:id", DeleteRule)
		fw.POST("rules/:id/toggle", ToggleRule)
	}
}
