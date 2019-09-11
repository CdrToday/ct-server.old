package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

func root(ctx iris.Context) {
	ctx.HTML("hello, wolrd")
}

type UserAPI struct {
	db *gorm.DB
}

func (u *UserAPI) mail(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	_uuid := uuid.NewV4().String()

	if sendMail(mail, _uuid) && rSet(mail, _uuid) {
		ctx.JSON(iris.Map{
			"msg": "ok",
		})
	}
}

/// @route: "/:mail/verify"
type VerifyBody struct {
	Code string `json: "code"`
}

func (u *UserAPI) verify(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body VerifyBody
	ctx.ReadJSON(&body)

	if rGet(mail) == body.Code {
		var user User

		u.db.FirstOrCreate(&user, User{Mail: mail})
		ctx.JSON(iris.Map{
			"msg": "ok",
			"data": iris.Map{
				"name": user.Name,
				"mail": user.Mail,
			},
		})

		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}
