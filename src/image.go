package main

import (
	"bytes"
	"encoding/json"
	"github.com/kataras/iris"
	"net/http"
)

type UploadBody struct {
	Image string `json:image`
}

func (u *UserAPI) upload(ctx iris.Context) {
	var body UploadBody
	ctx.ReadJSON(&body)

	pb := proxy(body)

	if pb.Msg == "err" {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	ctx.JSON(iris.Map{
		"msg":   "ok",
		"cover": pb.Msg,
	})
}

type ProxyBody struct {
	Msg string `json:msg`
}

func proxy(body UploadBody) ProxyBody {
	value, _ := json.Marshal(body)
	resp, _ := http.Post(
		"http://localhost:7070/upload",
		"application/json",
		bytes.NewBuffer(value),
	)

	var pb ProxyBody
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&pb)

	return pb
}
