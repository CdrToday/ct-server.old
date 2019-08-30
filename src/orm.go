package main

import (
	// "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

func main() {
	o := orm()
	user := User{
		Mail: "udtrokia@163.com",
	}
	o.First(&user)
	o.Model(&user).Update("articles", []string{"gorgor"})
}

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
	Name     string         `gorm:"unique"`
	Mail     string         `gorm:"unique"`
	Articles pq.StringArray `gorm:"type:varchar(100)[];"`
}

/// article
type Article struct {
	Id        string `gorm: "unique;"`
	Title     string
	Content   string
	Timestamp int64
}
