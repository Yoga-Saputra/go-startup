package config

import (
	"log"
	"startup/app/users"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDB() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// db_mysql_startup
	dsn := "root:P@ssW0rd@tcp(localhost:3306)/golang_startup?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&users.User{})
	return db
}
