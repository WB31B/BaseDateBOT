package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	ID      int    `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Message string `json:"user_message"`
}

func Connect() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:pG2r4hack@localhost:5432/postgres"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&User{})

	db.Create(&User{})

	DB = db
}
