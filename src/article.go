package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"strconv"
)

type PostAPI struct {
	db *gorm.DB
}

func (a *PostAPI) mail(ctx iris.Context) {
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

	posts := []Post{}
	var _arr []string = user.Posts
	a.db.Where("id IN (?)", _arr).Order("timestamp").Limit(limit).Offset(_page * limit).Find(&posts)

	ctx.JSON(iris.Map{
		"posts": posts,
	})
}

func (a *PostAPI) user(ctx iris.Context) {
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

	posts := []Post{}
	var _arr []string = user.Posts
	a.db.Where("id IN (?)", _arr).Order("timestamp").Limit(limit).Offset(_page * limit).Find(&posts)

	ctx.JSON(iris.Map{
		"posts": posts,
	})
}

func (a *PostAPI) spec(ctx iris.Context) {
	id := ctx.Params().Get("id")

	var post Post
	a.db.Where("id = ?", id).First(&post)

	ctx.JSON(iris.Map{
		"data": post,
	})
}
