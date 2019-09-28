package main

import (
	"encoding/base64"
	"github.com/kataras/iris"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"os"
)

type UploadBody struct {
	Image string `json:image`
}

func (u *UserAPI) upload(ctx iris.Context) {
	var body UploadBody
	ctx.ReadJSON(&body)

	key := proxy(body)

	ctx.JSON(iris.Map{
		"msg":   "ok",
		"image": key,
	})
}

func proxy(body UploadBody) string {
	t := conf()

	key, localFile := saveImage(body.Image)
	putPolicy := storage.PutPolicy{
		Scope: t.Get("qiniu.bucket").(string),
	}

	mac := qbox.NewMac(
		t.Get("qiniu.ak").(string),
		t.Get("qiniu.sk").(string),
	)

	upToken := putPolicy.UploadToken(mac)
	formUploader := storage.NewFormUploader(&storage.Config{})
	ret := storage.PutRet{}

	formUploader.PutFile(
		context.Background(),
		&ret,
		upToken,
		key,
		localFile,
		&storage.PutExtra{},
	)
	os.Remove(localFile)

	return ret.Key
}

func saveImage(b64 string) (string, string) {
	key := uuid.NewV4().String()
	home, _ := os.UserHomeDir()
	path := home + "/tmp/cache" + key + ".png"

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

func changeImageName(key string, dkey string) {
	t := conf()

	mac := qbox.NewMac(
		t.Get("qiniu.ak").(string),
		t.Get("qiniu.sk").(string),
	)

	bucket := t.Get("qiniu.bucket").(string)
	bucketManager := storage.NewBucketManager(mac, &storage.Config{})

	force := true
	bucketManager.Copy(bucket, key, bucket, dkey, force)
}
