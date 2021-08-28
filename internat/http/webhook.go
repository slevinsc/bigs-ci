package http

import (
	core2 "bigs-ci/internat/core"
	"bigs-ci/internat/dict"
	models2 "bigs-ci/internat/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"bigs-ci/lib/logger"
)

func PushEvent(c *gin.Context) {
	var req dict.PushEventPayload
	if err := c.ShouldBind(&req); err != nil {
		logger.Error("argument error", logger.Any("req", req), logger.Err(err))
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	logger.DebugAsJson(req)
	tGiblab := &models2.TGitlab{
		GitlabID:     models2.GetShortUUID(),
		ObjectKind:   req.ObjectKind,
		Before:       req.Before,
		After:        req.After,
		Ref:          req.Ref,
		CheckoutSHA:  req.CheckoutSHA,
		UserID:       req.UserID,
		UserName:     req.UserName,
		UserUsername: req.UserUsername,
		UserEmail:    req.UserEmail,
		UserAvatar:   req.UserAvatar,
		ProjectID:         req.ProjectID,
		TotalCommitsCount: req.TotalCommitsCount,
		ServiceName:       req.Project.Name,
	}
	project, err := json.Marshal(req.Project)
	if err != nil {
		logger.Error("toJsonProjectError", logger.Err(err), logger.Any("req", req))
		c.JSON(500, err)
		return
	}
	tGiblab.Project = string(project)
	repos, err := json.Marshal(req.Repository)
	if err != nil {
		logger.Error("toJsonReposError", logger.Err(err), logger.Any("req", req))
		c.JSON(500, err)
		return
	}
	tGiblab.Repository = string(repos)
	commits, err := json.Marshal(req.Commits)
	if err != nil {
		logger.Error("toJsonCommitError", logger.Err(err), logger.Any("req", req))
		c.JSON(500, err)
		return
	}
	tGiblab.Commits = string(commits)
	if err := tGiblab.Create(); err != nil {
		logger.Error("createGitlabError", logger.Err(err), logger.Any("req", req))
		c.JSON(500, err)
		return
	}
	var branch string
	if req.Ref != "" {
		s := strings.Split(req.Ref, "/")
		branch = s[len(s)-1]
	}

	go core2.EventHandle(&core2.EventData{
		ServiceName: req.Project.Name,
		Branch:      branch,
		CommitID:    req.CheckoutSHA,
		GitSshUrl:   req.Project.GitSSSHURL,
		GitHttpUrl:  req.Project.GitHTTPURL,
	})
	c.JSON(200, nil)
}
