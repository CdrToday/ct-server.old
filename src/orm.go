package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open("postgres", "dbname=cdr_today sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()
}
