package storage

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/clusterpedia-io/clusterpedia/pkg/storage"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jinzhu/configor"
	"gopkg.in/natefinch/lumberjack.v2"
	gmysql "gorm.io/driver/mysql"
	gpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8s.io/klog/v2"
)

const (
	StorageName = "sample-storage-layer"

	defaultLogFileName = "/var/log/clusterpedia/sample-storage-layer.log"
)

func RegisterStorageLayer() {
	storage.RegisterStorageFactoryFunc(StorageName, NewStorageFactory)
	klog.Infof("Successful register storage :%s", StorageName)
}

func NewStorageFactory(configPath string) (storage.StorageFactory, error) {
	if configPath == "" {
		return nil, fmt.Errorf("configPath should not be empty")
	}

	cfg := &Config{}
	if err := configor.Load(cfg, configPath); err != nil {
		return nil, err
	}

	var dialector gorm.Dialector
	switch cfg.Type {
	case "mysql":
		mysqlConfig, err := cfg.genMySQLConfig()
		if err != nil {
			return nil, err
		}

		connector, err := mysql.NewConnector(mysqlConfig)
		if err != nil {
			return nil, err
		}

		dialector = gmysql.New(gmysql.Config{Conn: sql.OpenDB(connector)})
	case "postgres":
		pgconfig, err := cfg.genPostgresConfig()
		if err != nil {
			return nil, err
		}

		dialector = gpostgres.New(gpostgres.Config{Conn: stdlib.OpenDB(*pgconfig)})
	default:
		return nil, fmt.Errorf("not support storage type: %s", cfg.Type)
	}

	logger, err := newLogger(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true, Logger: logger})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	connPool, err := cfg.getConnPoolConfig()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(connPool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(connPool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(connPool.ConnMaxLifetime)

	if err := db.AutoMigrate(&Resource{}); err != nil {
		return nil, err
	}

	return &StorageFactory{db}, nil
}

func newLogger(cfg *Config) (logger.Interface, error) {
	if cfg.Log == nil {
		return logger.Discard, nil
	}

	loggerConfig, err := cfg.LoggerConfig()
	if err != nil {
		return nil, err
	}

	var logWriter io.Writer
	if cfg.Log.Stdout {
		logWriter = os.Stdout
	} else {
		lumberjackLogger := cfg.Log.Logger
		if lumberjackLogger == nil {
			lumberjackLogger = &lumberjack.Logger{
				Filename:   defaultLogFileName,
				MaxSize:    100, // megabytes
				MaxBackups: 1,
			}
		} else if lumberjackLogger.Filename == "" {
			lumberjackLogger.Filename = defaultLogFileName
		}
		logWriter = lumberjackLogger
	}

	return logger.New(log.New(logWriter, "", log.LstdFlags), loggerConfig), nil
}
