package main

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"time"
)

type ReportBody struct {
	Type    string `json:"type"`
	Task    string `json:"task"`
	Content string `json:"content"`
}

func (u *UserAPI) report(ctx iris.Context) {
	ident := ctx.Params().Get("mail")

	var body ReportBody
	ctx.ReadJSON(&body)
	_uuid := uuid.NewV4().String()

	u.db.Create(&Report{
		Id:        _uuid,
		Type:      body.Type,
		Task:      body.Task,
		From:      ident,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	})

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}
