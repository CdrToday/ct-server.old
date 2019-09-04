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
	var user User
	a.db.Where("mail = ?", mail).Find(&user)

	articles := []Article{}
	var _arr []string = user.Articles
	a.db.Where("id IN (?)", _arr).Find(&articles)

	ctx.JSON(iris.Map{
		"articles": articles,
	})
}
