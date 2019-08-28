package main

import (
	"github.com/kataras/iris"
)

type Article struct{}

func (a *Article) articles(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	ctx.JSON(iris.Map{
		"msg":  "ok",
		"mail": mail,
	})
}
