package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"math/rand"
	"strconv"
)

type CommunityAPI struct {
	db *gorm.DB
}

// @route: POST "/u/:mail/c/create"
type CreateBody struct {
	Name string `json:name`
}

func genId(db *gorm.DB) string {
	id := rand.Intn(90000) + 10000
	var community Community
	if err := db.Where("id = ?", id).Find(&community).Error; err != nil {
		return strconv.Itoa(id)
	} else {
		return genId(db)
	}
}

func (c *CommunityAPI) create(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	_id := genId(c.db)

	var body CreateBody
	ctx.ReadJSON(&body)

	if err := c.db.Create(
		&Community{
			Id:         _id,
			Name:       body.Name,
			Owner:      mail,
			Avatar:     "",
			Members:    []string{mail},
			Applicants: []string{},
		},
	).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

type JoinBody struct {
	Id string `json:id`
}

// @route: POST "/u/:mail/c/join"
func (c *CommunityAPI) join(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body JoinBody
	ctx.ReadJSON(&body)

	var community Community
	c.db.Where("id = ?", body.Id).Find(&community)

	members := community.Members
	if b := contains(members, mail); b == true {
		ctx.JSON(iris.Map{
			"msg": "joined",
		})
		return
	}

	_members := append(community.Members, mail)
	if err := c.db.Model(&community).Where(
		"id = ?", body.Id,
	).Error; err == nil {
		if community.Id == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		c.db.Model(&community).Where(
			"id = ?", body.Id,
		).Update(
			"members", _members,
		)

		ctx.JSON(iris.Map{
			"msg": "ok",
		})
		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}

// updateCommunityName
type UpdateCommunityNameBody struct {
	Name string `json:name`
	Id   string `json:id`
}

func (u *CommunityAPI) updateCommunityName(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body UpdateCommunityNameBody
	ctx.ReadJSON(&body)

	var community Community
	if err := u.db.Model(&community).Where(
		"id = ?", body.Id,
	).Find(&community).Error; err == nil {
		if community.Owner != mail {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
		u.db.Model(&community).Where("id = ?", body.Id).Update("name", body.Name)
		u.db.Where("id = ?", body.Id).Select("name").Find(&community)

		ctx.JSON(iris.Map{
			"msg": "ok",
			"data": iris.Map{
				"name": community.Name,
			},
		})

		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}

// updateCommunityId
// type UpdateCommunityIdBody struct {
// 	targetId string `json:targetId`
// 	Id       string `json:id`
// }

// func (u *CommunityAPI) updateCommunityId(ctx iris.Context) {
// 	mail := ctx.Params().Get("mail")
// 	var body UpdateCommunityIdBody
// 	ctx.ReadJSON(&body)
//
// 	var community Community
// 	if err := u.db.Model(&community).Where(
// 		"id = ?", body.Id,
// 	).Find(&community).Error; err == nil {
// 		if community.Owner != mail {
// 			ctx.StatusCode(iris.StatusBadRequest)
// 			return
// 		}
//
// 		if err := u.db.Model(&community).Where(
// 			"id = ?", body.Id,
// 		).Update("id", body.targetId).Error; err == nil {
// 			print("changed")
// 			ctx.JSON(iris.Map{
// 				"msg": "ok",
// 			})
// 		}
//
// 		ctx.StatusCode(iris.StatusBadRequest)
// 	}
//
// 	ctx.StatusCode(iris.StatusBadRequest)
// }

// @route: GET "/u/:mail/c"
func (c *CommunityAPI) communities(ctx iris.Context) {
	mail := ctx.Params().Get("mail")

	var user User
	c.db.Where("mail = ?", mail).Find(&user)

	var communities []Community
	// var _arr []string = user.Communities
	if err := c.db.Where(
		"array_to_string(members, ',', '*') LIKE ?", "%"+mail+"%",
	).Select(
		"id, name, owner, avatar, members",
	).Find(&communities).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx.JSON(iris.Map{
		"communities": communities,
	})
}

// @route: GET "/u/:mail/c/:id"
func (c *CommunityAPI) members(ctx iris.Context) {
	id := ctx.Params().Get("id")

	var community Community
	c.db.Where("id = ?", id).Find(&community)

	var members []string = community.Members
	users := []User{}
	c.db.Where("mail in (?)", members).Select(
		"name, avatar, mail",
	).Find(&users)

	ctx.JSON(iris.Map{
		"members": users,
	})
}

// @route: GET "/u/:mail/c/:id/topics"
func (c *CommunityAPI) topics(ctx iris.Context) {
	id := ctx.Params().Get("id")

	var community Community
	c.db.Where("id = ?", id).Find(&community)

	var topics []Reddit
	var _topics []string = community.Topics
	c.db.Where("id in (?)", _topics).Order("timestamp desc").Find(&topics)

	ctx.JSON(iris.Map{
		"topics": topics,
	})
}

func (c *CommunityAPI) topicBatch(ctx iris.Context) {
	id := ctx.Params().Get("id")
	topic := ctx.Params().Get("topic")

	var community Community
	c.db.Where("id = ?", id).Find(&community)

	var reddits []Reddit
	c.db.Where("topic = ?", topic).Order("timestamp desc").Find(&reddits)

	ctx.JSON(iris.Map{
		"reddits": reddits,
	})
}

// @route: GET "/u/:mail/c/:id/quit"
func (c *CommunityAPI) quit(ctx iris.Context) {
	id := ctx.Params().Get("id")
	mail := ctx.Params().Get("mail")

	var community Community
	c.db.Where("id = ?", id).Find(&community)
	_ms := community.Members
	_ms = deleteStringFromArray(_ms, mail)

	if len(_ms) > 0 {
		c.db.Model(&Community{}).Where("id = ?", id).Update("members", _ms)
	} else {
		c.db.Model(&community).Delete(&Community{Id: id})
	}

	if mail == community.Owner && len(_ms) > 0 {
		c.db.Model(&Community{}).Where("id = ?", id).Update("owner", _ms[0])
	}

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}
