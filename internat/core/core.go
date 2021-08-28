package core

import (
	"bigs-ci/internat/dict"
	engine2 "bigs-ci/internat/engine"
	models2 "bigs-ci/internat/models"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"bigs-ci/config"
	"bigs-ci/lib/logger"
	"bigs-ci/lib/redis"
)

var TaskMap sync.Map

type Worker struct {
	L          *list.List
	Mutex      *sync.RWMutex
	TaskConfig *models2.TTaskConfig
}

func (w *Worker) serviceQueueHandle() {
	item := w.L.Front()
	de := item.Value.(*dockerEngine)

	err := eventHandle(de)
	update := make(map[string]interface{})
	if err != nil {
		logger.Error(fmt.Sprintf("%s deploy Error", de.ServiceName), logger.Err(err))
		update = map[string]interface{}{
			"build_state": dict.BuildStateFailure,
			"failure_at":  time.Now(),
		}

	} else {
		update = map[string]interface{}{
			"build_state": dict.BuildStateSuccessful,
			"success_at":  time.Now(),
		}
	}
	defer func() {
		w.Mutex.Lock()
		w.L.Remove(item)
		w.Mutex.Unlock()
		if w.L.Len() != 0 {
			w.serviceQueueHandle()
		} else {
			logger.Info("current task is empty")
		}
	}()
	tHistory := new(models2.TTaskHistory)
	r, err := tHistory.Update(de.HistoryID, update)
	if err != nil {
		logger.Error("updateHistoryError", logger.Err(err), logger.Any("update", update))
		return
	}
	if r == 0 {
		logger.Error("updateHistoryEmpty", logger.Any("update", update))
		return
	}
	return
}

func EventHandle(data *EventData) {
	serviceName := strings.ToLower(data.ServiceName)
	incr,err:=redis.Client.Cli.Incr(
		context.Background(),
		fmt.Sprintf("%s::%s", strings.ToLower(data.ServiceName), data.Branch)).Uint64()
	if err != nil {
		logger.Error("redisIncrError", logger.Err(err), logger.Any("data", data))
		return
	}
	tTaskConfig := new(models2.TTaskConfig)
	tc, err := tTaskConfig.Show(serviceName, data.Branch)
	if err != nil {
		logger.Error("queryTaskConfigError", logger.Err(err))
		return
	}
	var portNat []engine2.PortNat
	if err := json.Unmarshal([]byte(tc.Ports), &portNat); err != nil {
		logger.Error("parsePortNatError", logger.Err(err))
		return
	}

	h := &models2.TTaskHistory{
		TaskHistoryID: models2.GetShortUUID(),
		TaskID:        tc.TaskID,
		ServiceName:   serviceName,
		Version:       int64(incr),
		Branch:        data.Branch,
		GitUrl:        data.GitSshUrl,
		Commit:        data.CommitID,
	}
	if err := h.Create(); err != nil {
		logger.Error("createHistoryError", logger.Err(err))
		return
	}

	de := &dockerEngine{E: engine2.New(nil, config.Config.DockerAddr)}
	de.ServiceName = tc.ServiceName
	de.MainFileDir = tc.MainFileDir
	de.Language = dict.Language(tc.Language)
	de.Version = int64(incr)
	de.GitSshUrl = data.GitSshUrl
	de.Branch = data.Branch
	de.PortNat = portNat
	de.HistoryID = h.TaskHistoryID

	if task, ok := TaskMap.Load(tc.TaskName); ok {
		worker := task.(*Worker)
		worker.Mutex.Lock()
		worker.L.PushBack(de)
		worker.Mutex.Unlock()
		logger.Info(fmt.Sprintf("%s-%s-%d join queue,There are currently %d tasks queued", serviceName, de.Branch, de.Version, worker.L.Len()-1))
		if worker.L.Len() == 1 {
			go worker.serviceQueueHandle()
		}
	} else {
		worker := &Worker{
			L:          list.New(),
			TaskConfig: tTaskConfig,
			Mutex:      new(sync.RWMutex),
		}
		worker.Mutex.Lock()
		worker.L.PushBack(de)
		worker.Mutex.Unlock()
		logger.Info(fmt.Sprintf("%s-%s-%d init queue", serviceName, de.Branch, de.Version))
		TaskMap.Store(tc.TaskName, worker)
		go worker.serviceQueueHandle()
	}

}

func eventHandle(release Release) error {

	_, err := release.PullCode()
	if err != nil {
		logger.Error("pullCodeError", logger.Err(err))
		return err
	}

	if err := release.Build(); err != nil {
		logger.Error("buildCodeError", logger.Err(err))
		return err
	}

	if err := release.Deploy(); err != nil {
		logger.Error("deployCodeError", logger.Err(err))
		return err
	}
	return nil
}
