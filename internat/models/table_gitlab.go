package models

import (
	"bigs-ci/lib/db"
	"time"
)

type TGitlab struct {
	ID        uint64    `gorm:"primary_key;column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	GitlabID          string `gorm:"column:gitlab_id"`
	ObjectKind        string `gorm:"column:object_kind"`
	Before            string `gorm:"column:before"`
	After             string `gorm:"column:after"`
	Ref               string `gorm:"column:ref"`
	CheckoutSHA       string `gorm:"column:checkout_sha"`
	UserID            int64  `gorm:"column:user_id"`
	UserName          string `gorm:"column:user_name"`
	UserUsername      string `gorm:"column:user_username"`
	UserEmail         string `gorm:"column:user_email"`
	UserAvatar        string `gorm:"column:user_avatar"`
	ProjectID         int64  `gorm:"column:project_id"`
	Project           string `gorm:"column:project"`
	Repository        string `gorm:"column:repository"`
	Commits           string `gorm:"column:commits"`
	TotalCommitsCount int64  `gorm:"column:total_commits_count"`
	ServiceName       string `gorm:"column:service_name"`
}

func (t *TGitlab) TableName() string {
	return "gitlab"
}

func (t *TGitlab) Create() error {
	db := db.DB.MasterDB

	//if primary key is not existing
	err := db.Create(&t).Error
	return err
}
