package models

import (
	"bigs-ci/internat/dict"
	engine2 "bigs-ci/internat/engine"
	"bigs-ci/lib/logger"
	"testing"
)

func TestTTaskConfig_Create(t *testing.T) {
	ports := []engine2.PortNat{{}}
	j := new(dict.CreateTaskConfigReq)
	j.Ports = ports
	logger.DebugAsJson(j)
}
