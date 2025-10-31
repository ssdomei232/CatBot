package handler

import (
	"database/sql"
	"fmt"
	"time"

	"git.mmeiblog.cn/mei/CatBot/configs"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB() (*sql.DB, error) {
	config, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	// 验证配置是否完整
	if config.MYSQL.Host == "" || config.MYSQL.User == "" || config.MYSQL.DBName == "" {
		return nil, fmt.Errorf("mysql configuration is incomplete: host=%s, user=%s, dbname=%s",
			config.MYSQL.Host, config.MYSQL.User, config.MYSQL.DBName)
	}

	// 构建连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MYSQL.User,
		config.MYSQL.Password,
		config.MYSQL.Host,
		config.MYSQL.DBName) // 确保数据库名称在这里

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 测试连接
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
