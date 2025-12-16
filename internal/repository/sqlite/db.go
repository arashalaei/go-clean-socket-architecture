package store

import (
	"fmt"

	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IStore interface {
	repository.PersonRepositroy
	repository.SchoolRepository
}

type sqlit struct {
	db *gorm.DB
}

func NewSqlite(dbPath string) (IStore, error) {
	db, err := gorm.Open(sqlite.Dialector{
		DSN: dbPath,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sqlite db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	return &sqlit{db}, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.School{},
	)
}
