package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"nav_service/config"
	"sync"
)

var db *gorm.DB

var once sync.Once

func GetDB() *gorm.DB {
	once.Do(func() {
		dsn := config.GetConfig().Mysql.Dsn + "?charset=utf8mb4&parseTime=True&loc=Local"
		d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db = d
	})
	return db
}
