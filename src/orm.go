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
	db.AutoMigrate(&Reddit{})
	db.AutoMigrate(&Community{})
	db.AutoMigrate(&Report{})

	return db
}

/// user
type User struct {
	Avatar string         `json:"avatar"`
	Mail   string         `gorm:"unique"json:"mail"`
	Name   string         `json:"name"`
	Posts  pq.StringArray `gorm:"type:varchar(100)[];"json:"posts"`
}

/// article
type Post struct {
	Id        string `gorm:"unique;primary_key"json:"id"`
	Author    string `json:"author"`
	Document  string `json:"document"`
	Timestamp int64  `json:"timestamp"`
}

/// article
type Reddit struct {
	Id        string `gorm:"unique;primary_key"json:"id"`
	Type      string `json:"type"`
	Topic     string `json:"topic"`
	Author    string `json:"author"`
	Document  string `json:"document"`
	Community string `json:"community"`
	Timestamp int64  `json:"timestamp"`
}

/// community
type Community struct {
	Id         string         `gorm:"unique;primary_key"json:"id"`
	Name       string         `json:"name"`
	Owner      string         `json:"owner"`
	Avatar     string         `json:"avatar"`
	Topics     pq.StringArray `gorm:"type:varchar(100)[]"json:"topics"`
	Members    pq.StringArray `gorm:"type:varchar(100)[];"json:"members"`
	Applicants pq.StringArray `gorm:"type:varchar(100)[];"json:"applicants"`
}

/// report
type Report struct {
	Id        string `gorm:"unique;primary_key"json:"id"`
	Type      string `json:"type"`
	Task      string `json:"task"`
	From      string `json:"from"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
