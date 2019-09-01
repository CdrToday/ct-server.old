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
	var _arr []string = user.Articles
	a.db.Where("id IN (?)", _arr).Find(&articles)

	ctx.JSON(iris.Map{
		"articles": articles,
	})
}
