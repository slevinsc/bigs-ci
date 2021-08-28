package http

import (
	core2 "bigs-ci/internat/core"
	"bigs-ci/internat/dict"
	models2 "bigs-ci/internat/models"
	"bigs-ci/lib/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTaskHistory(c *gin.Context) {
	var req dict.CreateTaskHistoryReq
	if err := c.ShouldBind(&req); err != nil {
		logger.Error("argument error", logger.Any("req", req), logger.Err(err))
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	tHistory := new(models2.TTaskHistory)
	h, err := tHistory.Show(req.HistoryID)
	if err != nil {
		logger.Error("queryTaskHistoryError", logger.Err(err))
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	data := core2.EventData{
		ServiceName: h.ServiceName,
		Branch:      h.Branch,
		CommitID:    h.Commit,
		GitSshUrl:   h.GitUrl,
	}
	go core2.EventHandle(&data)

	c.JSON(200, nil)
}
