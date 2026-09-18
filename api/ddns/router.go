package ddns

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.RouterGroup) {
	d := r.Group("ddns")
	{
		d.GET("tasks", GetTasks)
		d.POST("tasks", CreateTask)
		d.PUT("tasks/:id", UpdateTask)
		d.DELETE("tasks/:id", DeleteTask)
		d.POST("tasks/:id/run", RunTaskNow)
		d.GET("interfaces", GetInterfaces)
	}
}
