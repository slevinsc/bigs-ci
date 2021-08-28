package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupMiddlewares(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	//r.Use(handlers.SetLocal)
}

//SetupRouters setup routers for endpoint
func SetupRouters(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "ok")
	})
	v1.POST("/webhook", PushEvent)
	v1.POST("/task_config", CreateTaskConfig)
	v1.POST("/task_history", CreateTaskHistory)
	v1.StaticFS("/code_explorer", http.Dir("/user/local/docker/deploy"))
}
