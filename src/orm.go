package main

import (
	// "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

func orm() *gorm.DB {
	t := conf()
	db, err := gorm.Open("postgres", t.Get("pg.addr").(string))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Article{})

	return db
}

/// user
type User struct {
	Name     string
	Mail     string         `gorm:"unique"`
	Articles pq.StringArray `gorm:"type:varchar(100)[];"`
}

/// article
type Article struct {
	Id        string `gorm:"unique;primary_key"json:"id"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
