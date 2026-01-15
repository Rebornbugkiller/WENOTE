package repo

import (
	"wenote-backend/config"
	"wenote-backend/internal/model"
	"wenote-backend/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	cfg := config.GlobalConfig.Database

	var gormLogger gormlogger.Interface
	if config.GlobalConfig.Server.Mode == "debug" {
		gormLogger = gormlogger.Default.LogMode(gormlogger.Info)
	} else {
		gormLogger = gormlogger.Default.LogMode(gormlogger.Silent)
	}

	db, err := gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("自动迁移表结构失败: %w", err)
	}

	logger.Info("数据库连接成功")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Notebook{},
		&model.Note{},
		&model.Tag{},
		&model.NoteTag{},
		&model.AuditLog{},
	)
}

func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
