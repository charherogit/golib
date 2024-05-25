package models

import (
	"fmt"
	"golib/config"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	SqlMaxOpenConns = 20
	SqlMaxIdleConns = 10
	SqlMaxLifetime  = 8 * time.Hour // 8小时
)

// 数据库配置
type DBConfig struct {
	Addr     string
	UserName string
	Password string
	DBName   string
}

// 表示连接数据库的字符串
func (c DBConfig) ToString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.Addr, c.DBName)
}

func CreateSqlDB(db string) (*sqlx.DB, error) {
	sqlConfig := DBConfig{
		Addr:     config.C.MysqlAddr,
		UserName: config.C.MysqlName,
		Password: config.C.MysqlPassword,
		DBName:   db,
	}
	clientDB, err := sqlx.Connect("mysql", sqlConfig.ToString())
	if err != nil {
		return nil, err
	}

	if err := clientDB.Ping(); err != nil {
		fmt.Printf("Ping err: %s\n", err)
		return nil, err
	}

	// 返回秒 28800s
	var timeoutTime int64
	var name string
	if err = clientDB.QueryRow("show global variables like 'wait_timeout';").Scan(&name, &timeoutTime); err != nil {
		fmt.Printf("Get err: %s\n", err)
		return nil, err
	}

	clientDB.SetMaxOpenConns(SqlMaxOpenConns)
	clientDB.SetMaxIdleConns(SqlMaxIdleConns)

	// 秒转为纳秒，设置为超过1半时间为交互过去
	clientDB.SetConnMaxLifetime(time.Duration(timeoutTime/2) * time.Second)
	return clientDB, nil
}
