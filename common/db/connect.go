package db

import (
	"fmt"
	"log"
	"main/config"

	"github.com/jinzhu/gorm"
)

func BuildConnect() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		config.MysqlCon.Username, config.MysqlCon.Password, config.MysqlCon.Addr, config.MysqlCon.Port, config.MysqlCon.Db, "10s")
	//mysql connection
	dba, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Conn mysql database error: %v", err)
	}
	return dba
}
