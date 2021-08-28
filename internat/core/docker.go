package core

import (
	"bigs-ci/internat/dict"
	engine2 "bigs-ci/internat/engine"
	"errors"
	"fmt"
	"strings"

	"bigs-ci/config"
	"bigs-ci/lib/logger"
)

type dockerEngine struct {
	E                 *engine2.DockerEngine
	PullContainerName string
	ServiceName       string
	MainFileDir       string
	Language          dict.Language
	Version           int64
	GitSshUrl         string
	Branch            string
	PortNat           []engine2.PortNat
	HistoryID         string
}

func (e *dockerEngine) PullCode() (string, error) {
	e.PullContainerName = fmt.Sprintf("pull-dev-%s-%d", e.ServiceName, e.Version)
	cc := &engine2.ContainerConfig{
		ContainerName: e.PullContainerName,
		ImageName:     "alpine-pull-dev",
		ExtraHosts:    []string{"test.org:192.168.1.3"},
	}
	defer func() {
		e.E.Logs(cc.ContainerName)
	}()

	var cmd string
	switch e.Language {
	case dict.TaskLanguageGo:
		gitlabAddr := strings.Split(config.Config.GitlabAddr, ":")
		cmdPrefix := fmt.Sprintf("ssh-keyscan -t rsa -p %s %s", gitlabAddr[1], gitlabAddr[0])
		cmd = fmt.Sprintf("%s > /root/.ssh/known_hosts &&"+
			"git clone -b dev %s %s-%d", cmdPrefix, e.GitSshUrl, e.ServiceName, e.Version)
		cc.WorkingDir = fmt.Sprintf("/data/%s/src", e.Language.String())

	default:
		cmdPrefix := fmt.Sprintf("ssh-keyscan -t rsa -p 22 test.org")
		cmd = fmt.Sprintf("%s > /root/.ssh/known_hosts &&"+
			"git clone -b dev %s %s-%d", cmdPrefix, e.GitSshUrl, e.ServiceName, e.Version)
		cc.WorkingDir = fmt.Sprintf("/data/%s", e.Language.String())
	}
	logger.Debug(fmt.Sprintf("%s-%s-%d", e.ServiceName, e.Branch, e.Version), logger.String("cmd", cmd))

	cc.Binds = []string{fmt.Sprintf("/usr/local/docker/deploy/%s:/data/%s", e.Language.String(), e.Language.String())}
	cc.Cmd = cmd

	respID, err := e.E.Create(cc, false)
	if err != nil {
		logger.Error("createPullContainerError", logger.Err(err))
		return "", err
	}

	logger.Info("pull code successful", logger.String("respID", respID), logger.Any("cc", cc))
	return respID, nil
}

func (e *dockerEngine) Build() error {
	//根据服务名查询语言,有的语言不需要编译
	cc := &engine2.ContainerConfig{
		ContainerName: fmt.Sprintf("build-dev-%s-%d", e.ServiceName, e.Version),
		VolumesFrom:   []string{e.PullContainerName},
	}
	var cmd string
	switch e.Language {
	case dict.TaskLanguageGo:
		cc.ImageName = "centos7-go-dev"
		cmd = fmt.Sprintf("rm -rf /data/go/src/%s &&ln -s /data/go/src/%s-%d /data/go/src/%s &&", e.ServiceName, e.ServiceName, e.Version, e.ServiceName)
		if e.MainFileDir != "" {
			cmd += fmt.Sprintf("cd /data/go/src/%s/%s &&", e.ServiceName, e.MainFileDir)
		} else {
			cmd += fmt.Sprintf("/data/go/src/%s &&", e.ServiceName)
		}

		cmd += fmt.Sprintf("pwd&&ls&&echo $GOPATH&&echo $GOROOT&&go build -o %s .", e.ServiceName)
	case dict.TaskLanguagePHP:
		cc.ImageName = "centos7-php-dev"
		cc.Env = []string{"COMPOSER_HOME=/data/php"}
		cmd = fmt.Sprintf("rm -rf /data/%s/%s &&cp -r /data/%s/%s-%d /data/%s/%s &&", e.Language.String(), e.ServiceName, e.Language.String(), e.ServiceName, e.Version, e.Language.String(), e.ServiceName)
		cmd += fmt.Sprintf("cd  /data/%s/%s &&composer update&&composer install&&cp /root/php/.env .&&php artisan key:generate", e.Language.String(), e.ServiceName)
	default:
		cc.ImageName = "centos7-go-dev"
		cmd = fmt.Sprintf("rm -rf /data/%s/%s &&cp -r /data/%s/%s-%d /data/%s/%s &&", e.Language.String(), e.ServiceName, e.Language.String(), e.ServiceName, e.Version, e.Language.String(), e.ServiceName)
		cmd += fmt.Sprintf("ls  /data/%s/%s/", e.Language.String(), e.ServiceName)
	}

	cc.Cmd = cmd
	logger.Debug(fmt.Sprintf("%s", e.ServiceName), logger.String("cmd", cmd))
	defer func() {

		e.E.Logs(cc.ContainerName)

		if err := e.E.Cleaner(cc.ContainerName); err != nil {
			return
		}

	}()
	respID, err := e.E.Create(cc, false)
	if err != nil {
		logger.Error("createBuildContainerError", logger.Err(err))
		return err
	}
	logger.Info("build code successful", logger.Any("cc", cc), logger.String("respID", respID))
	return nil
}

func (e *dockerEngine) Deploy() error {
	defer func() {
		if err := e.E.Cleaner(e.PullContainerName); err != nil {
			return
		}
	}()
	var containerName string

	switch e.Language {
	case dict.TaskLanguageGo:

		containerName = fmt.Sprintf("run-dev-%s", e.ServiceName)
	case dict.TaskLanguagePHP:
		containerName = fmt.Sprintf("run-dev-%s", e.ServiceName)
	case dict.TaskLanguageH5:
		return nil
	default:
		return errors.New("unknown Language")
	}

	containers, err := e.E.List(containerName, nil)
	if err != nil {
		logger.Error("queryContainerError", logger.Err(err))
		return err
	}
	if len(containers) != 0 {
		if containers[0].State != "running" && e.Language != dict.TaskLanguageH5 {
			err := e.E.Remove(containers[0])
			if err != nil {
				logger.Error("removeContainerError", logger.Err(err))
				return err
			}
		}

	}
	env := []string{"RUNMODE=dev"}

	cc := new(engine2.ContainerConfig)
	var binds []string
	switch e.Language {
	case dict.TaskLanguageGo:
		cc.ImageName = "centos7-go-dev"
		cc.VolumesFrom = []string{e.PullContainerName}
		cc.PortBinds = e.PortNat
		cc.Cmd = fmt.Sprintf("./%s", e.ServiceName)
		if e.MainFileDir != "" {
			cc.WorkingDir = fmt.Sprintf("/data/go/src/%s-%d/%s", e.ServiceName, e.Version, e.MainFileDir)
		} else {
			cc.WorkingDir = fmt.Sprintf("/data/go/src/%s-%d", e.ServiceName, e.Version)
		}
		env = append(env, "ETCD_ADDR=192.168.1.55:2379", "CONSUL_ADDR=192.168.1.56:8500")
	case dict.TaskLanguagePHP:
		cc.ImageName = "fpm-"
		cc.PortBinds = e.PortNat
		cc.Cmd = "chown -R www-data:www-data /var/www/html&&php-fpm"
		binds = append(binds, fmt.Sprintf("/usr/local/docker/deploy/%s/%s:/var/www/html", e.Language.String(), e.ServiceName))
		binds = append(binds, fmt.Sprintf("/usr/local/docker/deploy/%s/.composer:/var/www/html/.composer", e.Language.String()))
		env = append(env, "COMPOSER_HOME=/var/www/html")

	default:
		return errors.New("unknown Language")
	}
	binds = append(binds, "/etc/localtime:/etc/localtime")
	cc.Binds = binds
	cc.ContainerName = containerName
	cc.Restart = "always"
	cc.Env = env

	respID, err := e.E.Create(cc, true)
	if err != nil {
		logger.Error("createRunContainerError", logger.Err(err))
		return err
	}
	logger.Info("run code successful", logger.Any("cc", cc), logger.String("respID", respID))
	return nil
}
