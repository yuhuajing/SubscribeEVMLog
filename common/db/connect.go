package db

import (
	"fmt"
	"main/config"

	"github.com/jinzhu/gorm"
)

func Buildconnect() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		config.MysqlCon.Username, config.MysqlCon.Password, config.MysqlCon.Addr, config.MysqlCon.Port, config.MysqlCon.Db, "10s")
	//mysql connection
	dba, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Connect error:%s\n", err)
	}
	return dba
}
