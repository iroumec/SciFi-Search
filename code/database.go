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

type ContentType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"` // "movie", "book", etc.
}

type Work struct {
	ID            uint `gorm:"primaryKey"`
	Title         string
	ContentTypeID uint
	ContentType   ContentType `gorm:"foreignKey:ContentTypeID"`
}

type Rating struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	WorkID uint
	Work   Work `gorm:"foreignKey:WorkID"`
	Score  int
}

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../database/app.db"), &gorm.Config{})
	if err != nil {
		panic("Error al conectar a la base de datos.")
	}

	db.AutoMigrate(&User{}, &ContentType{}, &Work{}, &Rating{})

	return db
}
