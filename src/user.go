package main

import (
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

/// @route: "/:mail/publish"
type PublishBody struct {
	Title   string `json: "title"`
	Cover   string `json: "cover"`
	Content string `json: "content"`
}

func (u *UserAPI) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body PublishBody
	ctx.ReadJSON(&body)
	_uuid := uuid.NewV4().String()

	article := Article{
		Id:        _uuid,
		Title:     body.Title,
		Cover:     body.Cover,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	}

	var user User
	u.db.Where("mail = ?", mail).Find(&user)
	u.db.Create(&article)

	_articles := append(user.Articles, _uuid)
	u.db.Model(&user).Where("mail = ?", mail).Update("articles", _articles)

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

/// @route: "/{mail: string}/update/name"
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
		u.db.Where("name = ?", body.Name).Find(&user)

		ctx.JSON(iris.Map{
			"msg": "ok",
			"data": iris.Map{
				"mail": user.Mail,
				"name": user.Name,
			},
		})

		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}

/// @route: "/{mail: string}/article/update"
type UpdateArticleBody struct {
	Id      string `json:id`
	Title   string `json:title`
	Cover   string `json:cover`
	Content string `json:content`
}

func (u *UserAPI) updateArticle(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body UpdateArticleBody
	ctx.ReadJSON(&body)

	var user User
	var article Article

	u.db.Where("mail = ?", mail).Find(&user)
	var _arr []string = user.Articles
	if err := u.db.Where("id IN (?)", _arr).Find(&article).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	u.db.Model(&article).Where("id = ?", body.Id).Updates(map[string]interface{}{
		"title":   body.Title,
		"cover":   body.Cover,
		"content": body.Content,
	})

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

/// @route: "/{mail: string}/article/delete"
type DeleteArticleBody struct {
	Id string `json:id`
}

func (u *UserAPI) deleteArticle(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body DeleteArticleBody
	ctx.ReadJSON(&body)

	var user User
	article := Article{
		Id: body.Id,
	}

	u.db.Delete(&article)
	u.db.Where("mail = ?", mail).Find(&user)

	var index int
	for i, b := range user.Articles {
		if b == body.Id {
			index = i
		}
	}

	_arr := user.Articles
	_arr[index] = _arr[len(_arr)-1]
	_arr = _arr[:len(_arr)-1]

	u.db.Model(&user).Where("mail = ?", mail).Update("articles", _arr)
	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}
