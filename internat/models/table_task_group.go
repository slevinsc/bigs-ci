package models

import (
	"bigs-ci/lib/db"
	"time"
)

type TTaskConfigGroup struct {
	ID        uint64    `gorm:"primary_key;column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	GroupID   string `gorm:"column:group_id"`
	GroupName string `gorm:"column:group_name"`
}

func (t *TTaskConfigGroup) TableName() string {
	return "task_group_config"
}

func (t *TTaskConfigGroup) Create() error {
	db := db.DB.MasterDB

	//if primary key is not existing
	err := db.Create(&t).Error
	return err
}
