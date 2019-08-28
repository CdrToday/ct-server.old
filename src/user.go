package main

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

func root(ctx iris.Context) {
	ctx.HTML("hello, wolrd")
}

type User struct{}

func (u *User) sendCode(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	_uuid := uuid.NewV4()

	ctx.JSON(iris.Map{
		"msg":  "ok",
		"mail": mail,
		"uuid": _uuid,
	})
}

func (u *User) verify(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	ctx.JSON(iris.Map{
		"msg":  "ok",
		"mail": mail,
	})
}

func (u *User) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	ctx.JSON(iris.Map{
		"msg":  "ok",
		"mail": mail,
	})
}
