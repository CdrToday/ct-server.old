package main

import (
	"bytes"
	"encoding/json"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
)

type UploadBody struct {
	Image string `json:image`
}

func (u *UserAPI) upload(ctx iris.Context) {
	var body UploadBody
	ctx.ReadJSON(&body)

	pb := proxy(body)

	if pb.Message == "err" {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	ctx.JSON(iris.Map{
		"msg": "ok",
		"url": pb.Message,
	})
}

type ProxyBody struct {
	Message string `json:msg`
}

func proxy(body UploadBody) ProxyBody {
	value, _ := json.Marshal(body)
	resp, _ := http.Post(
		"http://localhost:7070/upload",
		"application/json",
		bytes.NewBuffer(value),
	)

	var pb ProxyBody
	_body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(_body, &pb)

	return pb
}
