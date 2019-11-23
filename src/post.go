package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"strconv"
)

type PostAPI struct {
	db *gorm.DB
}

func (a *PostAPI) posts(ctx iris.Context) {
	limit := 10
	id := ctx.Params().Get("ident")
	page := ctx.URLParamDefault("p", "0")
	_page, _ := strconv.Atoi(page)
	community := ctx.URLParamDefault("c", "95146")

	var posts []Reddit
	a.db.Where("author = ? AND community = ?", id, community).Order(
		"timestamp desc",
	).Limit(limit).Offset(_page * limit).Find(&posts)

	ctx.JSON(iris.Map{
		"posts": posts,
	})
}
