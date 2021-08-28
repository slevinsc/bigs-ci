package models

import "bigs-ci/lib/db"

func GetShortUUID() string {
	db := db.DB.MasterDB

	var uuid string
	db.Unscoped().Raw("SELECT UUID_SHORT()").Row().Scan(&uuid)
	return uuid
}
