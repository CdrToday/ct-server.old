package main

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"time"
)

// publish
type PublishBody struct {
	Document string `json:document`
}

func (u *UserAPI) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body PublishBody
	ctx.ReadJSON(&body)
	_uuid := uuid.NewV4().String()

	post := Post{
		Id:        _uuid,
		Document:  body.Document,
		Timestamp: time.Now().Unix(),
	}

	var user User
	u.db.Where("mail = ?", mail).Find(&user)
	u.db.Create(&post)

	_posts := append(user.Posts, _uuid)
	u.db.Model(&user).Where("mail = ?", mail).Update("posts", _posts)

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

/// updatePost
type UpdatePostBody struct {
	Document string `json:document`
}

func (u *UserAPI) updatePost(ctx iris.Context) {
	id := ctx.Params().Get("id")
	mail := ctx.Params().Get("mail")

	var body UpdatePostBody
	ctx.ReadJSON(&body)

	var user User
	post := Post{
		Id: id,
	}

	u.db.Where("mail = ?", mail).Find(&user)
	var _arr []string = user.Posts
	if err := u.db.Where("id IN (?)", _arr).Find(&post).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := u.db.Model(&post).Where("id = ?", id).Updates(map[string]interface{}{
		"document": body.Document,
	}).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

// delete post
func (u *UserAPI) deletePost(ctx iris.Context) {
	id := ctx.Params().Get("id")
	mail := ctx.Params().Get("mail")

	// delete post in user
	var user User
	u.db.Where("mail = ?", mail).Find(&user)

	index := 0
	for i, b := range user.Posts {
		if b == id {
			index = i
		}
	}

	if index == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	_arr := user.Posts
	_arr[index] = _arr[len(_arr)-1]
	_arr = _arr[:len(_arr)-1]

	u.db.Model(&user).Where("mail = ?", mail).Update("posts", _arr)

	// delete post
	post := Post{
		Id: id,
	}

	u.db.Delete(&post)

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}
