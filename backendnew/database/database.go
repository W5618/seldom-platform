package database

import (
	"fmt"
	"seldom-platform/config"
	"seldom-platform/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var err error
	var dsn string

	switch cfg.Driver {
	case "sqlite3":
		dsn = cfg.Database
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	DB, err = gorm.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移数据库表
	err = autoMigrate()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return DB, nil
}

func autoMigrate() error {
	err := DB.AutoMigrate(
		&models.Project{},
		&models.Env{},
		&models.TestCase{},
		&models.TestCaseTemp{},
		&models.CaseResult{},
		&models.TestTask{},
		&models.TaskCaseRelevance{},
		&models.TaskReport{},
		&models.ReportDetails{},
		&models.Team{},
		&models.User{},
	).Error
	return err
}

func Close(db *gorm.DB) {
	if db != nil {
		db.Close()
	}
}

func GetDB() *gorm.DB {
	return DB
}