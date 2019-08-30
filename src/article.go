package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

type ArticleAPI struct {
	db *gorm.DB
}

func (a *ArticleAPI) articles(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	// get user
	user := User{
		Mail: mail,
	}
	a.db.First(&user)

	articles := []Article{}
	a.db.Find(&articles).Where("id IN (?)", user.Articles)

	ctx.JSON(iris.Map{
		"articles": articles,
	})
}
