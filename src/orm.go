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
	db.AutoMigrate(&Community{})

	return db
}

/// user
type User struct {
	Avatar      string
	Mail        string `gorm:"unique"`
	Name        string
	Posts       pq.StringArray `gorm:"type:varchar(100)[];"`
	Communities pq.StringArray `gorm:"type:varchar(100)[];"`
}

/// article
type Post struct {
	Id        string `gorm:"unique;primary_key"json:"id"`
	Document  string `json:"document"`
	Timestamp int64  `json:"timestamp"`
}

/// community
type Community struct {
	Id         string         `gorm:"unique;primary_key"json:"id"`
	Name       string         `json:"name"`
	Owner      string         `json:"owner"`
	Avatar     string         `json:"avatar"`
	Members    pq.StringArray `gorm:"type:varchar(100)[];"json:"members"`
	Applicants pq.StringArray `gorm:"type:varchar(100)[];"json:"applicants"`
}
