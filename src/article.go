package main

import (
	"github.com/kataras/iris"
)

type ArticleAPI struct{}

func (a *ArticleAPI) articles(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	ctx.JSON(iris.Map{
		"msg":  "ok",
		"mail": mail,
	})
}
