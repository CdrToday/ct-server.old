package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

type ArticleAPI struct {
	db *gorm.DB
}

func (a *ArticleAPI) mail(ctx iris.Context) {
	_mail := ctx.Params().Get("mail")

	// get user
	var user User
	a.db.Where("mail = ?", _mail).Find(&user)

	articles := []Article{}
	var _arr []string = user.Articles
	a.db.Where("id IN (?)", _arr).Find(&articles)

	ctx.JSON(iris.Map{
		"articles": articles,
	})
}

func (a *ArticleAPI) user(ctx iris.Context) {
	name := ctx.Params().Get("user")

	// get user
	var user User
	a.db.Where("name = ?", name).Find(&user)

	articles := []Article{}
	var _arr []string = user.Articles
	a.db.Where("id IN (?)", _arr).Find(&articles)

	ctx.JSON(iris.Map{
		"articles": articles,
	})
}

func (a *ArticleAPI) spec(ctx iris.Context) {
	id := ctx.Params().Get("id")

	var article Article
	a.db.Where("id = ?", id).First(&article)

	ctx.JSON(iris.Map{
		"data": article,
	})
}
