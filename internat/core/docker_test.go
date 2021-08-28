package core

import (
	"bigs-ci/config"
	engine2 "bigs-ci/internat/engine"
	"testing"
)

func init() {
	config.Init()
}

func TestDockerEngine_PullCode(t *testing.T) {
	de := dockerEngine{
		E:                 engine2.New(nil, "tcp://192.168.1.2:2375"),
		PullContainerName: "pull-test",
		ServiceName:       "test",
		MainFileDir:       "",
		Language:          1,
		Version:           1,
		GitSshUrl:         "ssh://git@test.org:1422/root/test.git",
		Branch:            "dev",
	}
	respID, err := de.PullCode()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(respID)
}
