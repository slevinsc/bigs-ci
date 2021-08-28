package models

import (
	"bigs-ci/lib/db"
	"time"
)

type TTaskConfig struct {
	ID        uint64    `gorm:"primary_key;column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	TaskID      string `gorm:"column:task_id"`
	TaskGroupID string `gorm:"column:task_group_id"`
	TaskName    string `gorm:"column:task_name"`
	ServiceName string `gorm:"column:service_name"`
	Language    uint8  `gorm:"column:language"`
	Ports       string `gorm:"column:ports"` // []portNAT json数组
	ImageName   string `gorm:"column:image_name"`
	Branch      string `gorm:"column:branch"`
	MainFileDir string `gorm:"column:main_file_dir"` // 存放main.go的路径
}

func (t *TTaskConfig) TableName() string {
	return "task_config"
}

func (t *TTaskConfig) Create() error {
	db := db.DB.MasterDB

	//if primary key is not existing
	err := db.Create(&t).Error
	return err
}

func (t *TTaskConfig) Show(serviceName, branch string) (*TTaskConfig, error) {
	db := db.DB.SlaveDB
	ret := new(TTaskConfig)
	err := db.Model(&t).Where("service_name=? and branch=?", serviceName, branch).Scan(ret).Error
	return ret, err
}
