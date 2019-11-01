package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
)

type RedditAPI struct {
	db *gorm.DB
}

///@route: GET "/u/:mail/c/:id/reddit"
func (r *RedditAPI) reddits(ctx iris.Context) {
	limit := 10
	id := ctx.Params().Get("id")
	page := ctx.URLParamDefault("p", "0")
	_page, _ := strconv.Atoi(page)

	var reddits []Reddit
	r.db.Where("community = ?", id).Order("timestamp desc").Limit(limit).Offset(_page * limit).Find(&reddits)

	ctx.JSON(iris.Map{
		"msg":     "ok",
		"reddits": reddits,
	})
}

///@route: POST "/u/:mail/reddit"
type RedditBody struct {
	Type      string `json:"type"`
	Community string `json:"community"`
	Document  string `json:"document"`
}

func (r *RedditAPI) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	var body RedditBody
	ctx.ReadJSON(&body)
	_uuid := uuid.NewV4().String()

	r.db.Create(&Reddit{
		Id:        _uuid,
		Type:      body.Type,
		Author:    mail,
		Community: body.Community,
		Document:  body.Document,
		Timestamp: time.Now().Unix(),
	})

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

/// @route: PUT "/u/:mail/r/:id"
type UpdateReddit struct {
	Document string `json:"document"`
}

func (r *RedditAPI) updateReddit(ctx iris.Context) {
	id := ctx.Params().Get("id")

	var body UpdateReddit
	ctx.ReadJSON(&body)

	r.db.Model(&Reddit{}).Where("id = ?", id).Updates(
		map[string]interface{}{"document": body.Document},
	)

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

func (r *RedditAPI) updateRedditTime(ctx iris.Context) {
	id := ctx.Params().Get("id")

	var body UpdateReddit
	ctx.ReadJSON(&body)

	r.db.Model(&Reddit{}).Where("id = ?", id).Update("timestamp", time.Now().Unix())

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

/// @route: DELETE "/u/:mail/r/:id"
func (r *RedditAPI) deleteReddit(ctx iris.Context) {
	id := ctx.Params().Get("id")
	r.db.Delete(&Reddit{Id: id})

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}
