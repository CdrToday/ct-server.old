package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

func orm() *gorm.DB {
	t := conf()
	db, _ := gorm.Open("postgres", t.Get("pg.addr").(string))

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})

	return db
}

/// user
type User struct {
	Avatar string
	Mail   string `gorm:"unique"`
	Name   string
	Posts  pq.StringArray `gorm:"type:varchar(100)[];"`
}

/// article
type Post struct {
	Id        string `gorm:"unique;primary_key"json:"id"`
	Document  string `json:"document"`
	Timestamp int64  `json:"timestamp"`
}
