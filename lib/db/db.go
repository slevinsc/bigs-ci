package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"bigs-ci/config"
	"bigs-ci/lib/logger"
)

type Store struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
	config   *config.Configure
}

type MySQLConfig struct {
	DSN     string `json:"dsn"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	DBName  string `json:"dbname"`
	MaxOpen int    `json:"max_open"`
	MaxIdle int    `json:"max_idle"`
}

var DB *Store

func Init() {
	DB = NewStore()
}

func NewStore() *Store {

	store := &Store{config: config.Config}
	store.NewMySQL()

	logger.Info("new bigs-ci mysql store", logger.Any(
		"MasterDB", config.Config.MysqlConfig.Master),
		logger.Any("SlaveDB", config.Config.MysqlConfig.Slave),
	)

	return store
}

func (s *Store) NewMySQL() {
	//?parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=True&timeout=5s&charset=utf8mb4&loc=Local",
		s.config.MysqlConfig.Master.User,
		s.config.MysqlConfig.Master.Pass,
		"tcp",
		s.config.MysqlConfig.Master.Host,
		s.config.MysqlConfig.Master.Port,
		s.config.MysqlConfig.Master.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mysqldb, err := db.DB()
	if err != nil {
		panic(err)
	}

	mysqldb.SetConnMaxLifetime(5 * time.Hour)
	mysqldb.SetMaxIdleConns(s.config.MysqlConfig.Master.MaxIdle)

}
