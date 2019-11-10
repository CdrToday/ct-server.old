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
	Topic     string `json:"topic"`
	Document  string `json:"document"`
	Community string `json:"community"`
}

func (r *RedditAPI) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	var body RedditBody
	ctx.ReadJSON(&body)
	_uuid := uuid.NewV4().String()

	var hasTopic bool
	var community Community
	if body.Topic != "" {
		r.db.Where("id = ?", body.Community).Find(&community)

		for _, t := range community.Topics {
			if t == body.Topic {
				hasTopic = true
			}
		}

		if !hasTopic {
			_topics := append(community.Topics, body.Topic)
			r.db.Model(&community).Where(
				"id = ?", body.Community,
			).Update(
				"topics", _topics,
			)

			r.db.Model(Reddit{}).Where(
				"id = ?", body.Topic,
			).Update("topic", body.Topic)
		}
	}

	r.db.Create(&Reddit{
		Id:        _uuid,
		Type:      body.Type,
		Topic:     body.Topic,
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
