package logcenter

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.RouterGroup) {
	l := r.Group("logcenter")
	{
		l.GET("sources", GetSources)
		l.GET("query", QueryLogs)
	}
}
