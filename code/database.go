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
	Unit          bool

	// Relación self-reference
	SagaID *uint
	Saga   *Work  `gorm:"foreignKey:SagaID"`
	Books  []Work `gorm:"foreignKey:SagaID"`
}

/*
Un Work puede ser saga o unidad. Si es saga, tendrá hijos; si es libro, puede apuntar a su saga

Si SagaID es NULL → es una saga.

Si SagaID apunta a otro Work → es un libro dentro de esa saga.

Books permite acceder desde la saga a todos los libros.

SagaID: es la FK que dice “esta obra pertenece a tal saga”.

Saga: es la referencia inversa a esa saga (el padre).

Books: es el slice de hijos que tiene esta obra (si es saga).

En otras palabras:

Si SagaID == nil → el registro es una saga (colección).

Podés acceder a sus hijos con db.Preload("Books").Find(&work)

Si SagaID != nil → el registro es un libro/unidad.

Podés acceder a la saga a la que pertenece con db.Preload("Saga").Find(&work)

*/

// Una colección solo puede estar compuesto por elementos de un mismo tipo.
// ¿Qué pasa con los videojuegos? Que tienen DLC, OSL...

type AssociatedWorks struct {
	ID     uint
	WorkID uint
	Work   Work `gorm:"foreignKey:WorkID"`
}

type Rating struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	WorkID uint
	Work   Work `gorm:"foreignKey:WorkID"`
	Score  int
	Review string
}

func example() {

	// Crear saga
	saga := Work{Title: "El Señor de los Anillos"}
	db.Create(&saga)

	// Crear libros asociados a la saga
	book1 := Work{Title: "La Comunidad del Anillo", SagaID: &saga.ID}
	book2 := Work{Title: "Las Dos Torres", SagaID: &saga.ID}
	book3 := Work{Title: "El Retorno del Rey", SagaID: &saga.ID}
	db.Create(&book1)
	db.Create(&book2)
	db.Create(&book3)

	var s Work
	db.Preload("Books").First(&s, saga.ID)
	// s.Books tendrá los 3 libros

	var b Work
	db.Preload("Saga").First(&b, book1.ID)
	// b.Saga.Title == "El Señor de los Anillos"

	/*
		First(&s, saga.ID) → busca el Work cuyo ID == saga.ID (es decir, la saga).

		Preload("Books") → además carga todos los Work que tengan SagaID == saga.ID.

		Resultado:

		s es la saga El Señor de los Anillos.

		s.Books contiene los tres libros hijos.
	*/

}

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../database/app.db"), &gorm.Config{})
	if err != nil {
		panic("Error al conectar a la base de datos.")
	}

	db.AutoMigrate(&User{}, &ContentType{}, &Work{}, &Rating{})

	return db
}
