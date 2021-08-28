package models

import (
	"bigs-ci/lib/db"
	"time"
)

type TTaskHistory struct {
	ID        uint64    `gorm:"primary_key;column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	TaskHistoryID string    `gorm:"column:task_history_id"`
	TaskID        string    `gorm:"column:task_id"`
	ServiceName   string    `gorm:"column:service_name"`
	Version       int64     `gorm:"column:version"` // 版本号
	BuildState    uint8     `gorm:"column:build_state"`
	SuccessAt     time.Time `gorm:"column:success_at"`
	FailureAt     time.Time `gorm:"column:failure_at"`
	Branch        string    `gorm:"column:branch"`
	GitUrl        string    `gorm:"column:git_url"`
	Commit        string    `gorm:"column:commit"`  // 提交ID
	RespID        string    `gorm:"column:resp_id"` //docker容器ID

}

func (t *TTaskHistory) TableName() string {
	return "task_history"
}

func (t *TTaskHistory) Create() error {
	db := db.DB.MasterDB

	//if primary key is not existing
	err := db.Create(&t).Error
	return err
}

func (t *TTaskHistory) Update(taskHistoryID string, update map[string]interface{}) (int64, error) {
	update["updated_at"] = time.Now()
	db := db.DB.MasterDB
	r := db.Model(&t).Where("task_history_id=?", taskHistoryID).Updates(update)
	return r.RowsAffected, r.Error
}

func (t *TTaskHistory) Show(ID string) (*TTaskHistory, error) {
	db := db.DB.SlaveDB
	ret := new(TTaskHistory)
	err := db.Model(&t).Where("task_history_id=?", ID).Scan(ret).Error
	return ret, err
}

func (t *TTaskHistory) List(serviceName, branch string) ([]*TTaskHistory, error) {
	db := db.DB.SlaveDB
	var ret []*TTaskHistory
	err := db.Model(&t).Where("service_name=? and branch=?", serviceName, branch).Scan(ret).Error
	return ret, err
}
