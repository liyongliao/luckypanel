package docker

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.RouterGroup) {
	d := r.Group("docker")
	{
		d.GET("status", GetStatus)
		d.GET("containers", ListContainers)
		d.POST("containers/:id/action", ContainerAction)
		d.GET("containers/:id/logs", GetLogs)
		d.GET("images", ListImages)
		d.GET("compose", ListCompose)
		d.POST("compose", SaveCompose)
		d.POST("compose/:name/action", ComposeAction)
	}
}
