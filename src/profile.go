package main

import (
	"github.com/kataras/iris"
)

// updateUser
type UpdateUserNameBody struct {
	Name string `json:name`
}

func (u *UserAPI) updateUserName(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body UpdateUserNameBody
	ctx.ReadJSON(&body)

	var user User
	if err := u.db.Where("name = ?", body.Name).Find(&user).Error; err != nil {
		u.db.Model(&user).Where("mail = ?", mail).Update("name", body.Name)
		u.db.Where("name = ?", body.Name).Select("name").Find(&user)

		ctx.JSON(iris.Map{
			"msg": "ok",
			"data": iris.Map{
				"name": user.Name,
			},
		})

		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}

// updateUser
type UpdateUserAvatarBody struct {
	Avatar string `json:avatar`
}

func (u *UserAPI) updateUserAvatar(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body UpdateUserAvatarBody
	ctx.ReadJSON(&body)

	u.db.Where("mail = ?", mail).Find(&User{}).Update("avatar", body.Avatar)

	ctx.JSON(iris.Map{
		"msg":    "ok",
		"avatar": body.Avatar,
	})
}
