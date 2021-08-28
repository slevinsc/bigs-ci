package http

import (
	"bigs-ci/internat/dict"
	models2 "bigs-ci/internat/models"
	"bigs-ci/lib/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTaskConfig(c *gin.Context) {
	var req dict.CreateTaskConfigReq
	if err := c.ShouldBind(&req); err != nil {
		logger.Error("argument error", logger.Any("req", req), logger.Err(err))
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	tConfig := &models2.TTaskConfig{
		TaskID:      models2.GetShortUUID(),
		TaskGroupID: "1",
		TaskName:    req.TaskName,
		ServiceName: req.ServiceName,
		Language:    req.Language,
		Branch:      req.Branch,
		MainFileDir: req.MainFileDir,
	}
	ports, err := json.Marshal(req.Ports)
	if err != nil {
		logger.Error("toJsonError", logger.Err(err), logger.Any("req", req))
		c.JSON(500, err)
		return
	}
	tConfig.Ports = string(ports)
	if err := tConfig.Create(); err != nil {
		logger.Error("createTaskConfigError",
			logger.Err(err),
			logger.Any("req", req))
		c.JSON(500, err)
		return
	}
	c.JSON(200, nil)
}
