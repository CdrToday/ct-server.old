package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"strconv"
)

type ArticleAPI struct {
	db *gorm.DB
}

func (a *ArticleAPI) mail(ctx iris.Context) {
	limit := 10
	_mail := ctx.Params().Get("mail")
	page := ctx.URLParamDefault("p", "0")
	_page, err := strconv.Atoi(page)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	// get user
	var user User
	a.db.Where("mail = ?", _mail).Find(&user)

	articles := []Article{}
	var _arr []string = user.Articles
	a.db.Where("id IN (?)", _arr).Order("timestamp").Limit(limit).Offset(_page * limit).Find(&articles)

	ctx.JSON(iris.Map{
		"articles": articles,
	})
}

func (a *ArticleAPI) user(ctx iris.Context) {
	limit := 10
	name := ctx.Params().Get("user")
	page := ctx.URLParamDefault("p", "0")
	_page, err := strconv.Atoi(page)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	// get user
	var user User
	a.db.Where("name = ?", name).Find(&user)

	articles := []Article{}
	var _arr []string = user.Articles
	a.db.Where("id IN (?)", _arr).Order("timestamp").Limit(limit).Offset(_page * limit).Find(&articles)

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
