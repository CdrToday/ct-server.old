package main

import (
	// "fmt"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

type CommunityAPI struct {
	db *gorm.DB
}

// @route: POST "/u/:mail/c/create"
type CreateBody struct {
	Id   string `json:id`
	Name string `json:name`
}

func (c *CommunityAPI) create(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body CreateBody
	ctx.ReadJSON(&body)

	var community Community
	if err := c.db.Where("id = ?", body.Id).Find(&community).Error; err == nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := c.db.Create(
		&Community{
			Id:         body.Id,
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

// @route: POST "/u/:mail/c/join"
func (c *CommunityAPI) join(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body CreateBody
	ctx.ReadJSON(&body)

	var community Community
	c.db.Where("id = ?", body.Id).Find(&community)

	// applicants := community.Applicants
	// _applicants := append(applicants, mail)
	// c.db.Model(&community).Where("id = ?", body.Id).Update("applicants", _applicants)

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
	).Update(
		"members", _members,
	).Error; err == nil {
		ctx.JSON(iris.Map{
			"msg": "ok",
		})
		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}

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
