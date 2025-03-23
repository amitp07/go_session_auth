package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup() *gorm.DB {
	dsn := "postgres://goSessionAuthU:goSessionAuthP@localhost:5432/go_session_auth"
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("Could not connect to DB..")
	}

	return db
}
