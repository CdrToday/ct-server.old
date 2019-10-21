package main

import (
	"encoding/base64"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type UploadBody struct {
	Image string `json:image`
}

func (u *UserAPI) upload(ctx iris.Context) {
	var body UploadBody
	ctx.ReadJSON(&body)

	key := proxy(body)

	if key == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx.JSON(iris.Map{
		"msg":   "ok",
		"image": key,
	})
}

func proxy(body UploadBody) string {
	t := conf()

	key, localFile := saveImage(body.Image)

	u, _ := url.Parse(t.Get("cos.url").(string))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  t.Get("cos.sid").(string),
			SecretKey: t.Get("cos.sk").(string),
		},
	})

	name := "images/" + key
	f := strings.NewReader("images")

	_, err := c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}

	_, err = c.Object.PutFromFile(context.Background(), name, localFile, nil)
	if err != nil {
		panic(err)
	}

	os.Remove(localFile)
	return key
}

func saveImage(b64 string) (string, string) {
	key := uuid.NewV4().String()
	home, _ := os.UserHomeDir()
	_dir := home + "/tmp/cache/"
	_, err := os.Stat(_dir)
	if err != nil {
		os.MkdirAll(_dir, os.ModePerm)
	}

	path := _dir + key + ".png"

	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	f, _ := os.Create(path)
	defer f.Close()

	f.Write(dec)
	f.Sync()

	return key, path
}
