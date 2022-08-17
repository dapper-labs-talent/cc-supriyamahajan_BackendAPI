package db

import (
	"cc-supriyamahajan_BackendAPI/models"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var DB *pg.DB
var dbError error

func Connect() *pg.DB {
	opt, err := pg.ParseURL("postgres://postgres:@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)
	if db == nil {
		log.Fatal("cannot connect to DB")
	}

	err = createSchema(db)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to Database!")
	DB = db
	return db
}

// createSchema creates database schema for User model.
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
