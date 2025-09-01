package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}

type Work struct {
	ID    uint `gorm:"primaryKey"`
	Title string
	Type  string
}

type Rating struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	WorkID uint
	Score  int
}

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../database/app.db"), &gorm.Config{})
	if err != nil {
		panic("Error al conectar a la base de datos.")
	}

	db.AutoMigrate(&User{})

	return db
}
