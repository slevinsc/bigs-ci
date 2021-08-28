package dict

import (
	engine2 "bigs-ci/internat/engine"
)

const (
	TaskLanguageGo  = 1
	TaskLanguageH5  = 2
	TaskLanguagePHP = 3
)

type Language uint8

func (l *Language) String() string {
	switch int(*l) {
	case TaskLanguageGo:
		return "go"
	case TaskLanguageH5:
		return "h5"
	case TaskLanguagePHP:
		return "php"
	default:
		return ""
	}
}

func LanguageToInt(key string) uint8 {
	lMap := map[string]uint8{
		"go":  TaskLanguageGo,
		"h5":  TaskLanguageH5,
		"php": TaskLanguagePHP,
	}
	return lMap[key]
}

const (
	BuildStateSuccessful = 1
	BuildStateFailure    = 2
)

type CreateTaskConfigReq struct {
	TaskName    string            `json:"task_name"`
	ServiceName string            `json:"service_name"`
	Language    uint8             `json:"language"`
	Ports       []engine2.PortNat `json:"ports"`
	Branch      string            `json:"branch"`
	MainFileDir string            `json:"main_file_dir"`
}
