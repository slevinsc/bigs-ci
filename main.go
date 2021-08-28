package main

import (
	http2 "bigs-ci/internat/http"
	"github.com/gin-gonic/gin"

	"bigs-ci/config"
	"bigs-ci/lib/db"
	"bigs-ci/lib/redis"
)

func main() {
	config.Init()
	db.Init()
	redis.Init()
	r := gin.Default()
	http2.SetupMiddlewares(r)
	http2.SetupRouters(r)
	r.Run("0.0.0.0:14880") // listen and serve on
}
