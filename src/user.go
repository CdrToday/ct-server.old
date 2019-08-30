package main

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

func root(ctx iris.Context) {
	ctx.HTML("hello, wolrd")
}

type UserAPI struct{}

func (u *UserAPI) sendCode(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	_uuid := uuid.NewV4().String()

	res := sendMail(mail, _uuid)
	_res := rSet(mail, _uuid)

	if res && _res {
		ctx.JSON(iris.Map{
			"msg": "ok",
		})
	}
}

type VerifyBody struct {
	Code string `json: "code"`
}

func (u *UserAPI) verify(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body VerifyBody
	err := ctx.ReadJSON(&body)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	res := rGet(mail)
	if res == body.Code {
		ctx.JSON(iris.Map{
			"msg": "ok",
		})
	}
}

func (u *UserAPI) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	ctx.JSON(iris.Map{
		"msg":  "ok",
		"mail": mail,
	})
}
