package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"time"
)

func root(ctx iris.Context) {
	ctx.HTML("hello, wolrd")
}

type UserAPI struct {
	db *gorm.DB
}

func (u *UserAPI) sendCode(ctx iris.Context) {
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
	fmt.Println(body)
	// pair uuid in redis
	if rGet(mail) == body.Code {
		user := User{
			Mail: mail,
		}

		u.db.FirstOrCreate(&user, User{Mail: mail})
		ctx.JSON(iris.Map{
			"msg": "ok",
		})

		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}

/// @route: "/:mail/publish"
type PublishBody struct {
	Title   string `json: "title"`
	Content string `json: "content"`
}

func (u *UserAPI) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body PublishBody
	ctx.ReadJSON(&body)
	_uuid := uuid.NewV4().String()

	user := User{
		Mail: mail,
	}

	article := Article{
		Id:        _uuid,
		Title:     body.Title,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	}

	// intt
	u.db.Find(&user)
	u.db.Create(&article)

	_articles := append(user.Articles, ArticleId{Id: _uuid})
	err := u.db.Model(&user).Update("articles", _articles)

	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}
