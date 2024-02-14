package db

import (
	"fmt"

	"github.com/zeabix-cloud-native/go-common/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormRead struct {
	*gorm.DB
}
type GormWrite struct {
	*gorm.DB
}

func NewConnectReadDB(cfg config.Config) (*GormRead, error) {
	psqlRead := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHostRead, cfg.DBUserRead, cfg.DBNameRead, cfg.DBPortRead, cfg.DBPasswordRead)
	dbread, dbRErr := gorm.Open(postgres.Open(psqlRead), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	return &GormRead{DB: dbread}, dbRErr
}

func NewConnectWriteDB(cfg config.Config) (*GormWrite, error) {
	psqlWrite := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHostWrite, cfg.DBUserWrite, cfg.DBNameWrite, cfg.DBPortWrite, cfg.DBPasswordWrite)
	dbwrite, dbWErr := gorm.Open(postgres.Open(psqlWrite), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	return &GormWrite{DB: dbwrite}, dbWErr
}
