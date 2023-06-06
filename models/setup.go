package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
const DSN =  "root:root@tcp(127.0.0.1:3306)/go_posts?charset=utf8mb4&parseTime=True&loc=Local"


func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Post{})

	DB = database
}
