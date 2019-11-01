package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Logger()
	app.Use(recover.New())
	app.Use(logger.New())

	// Middleware
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	// Methods
	_orm := orm()
	// _orm.LogMode(true)
	defer _orm.Close()

	user := UserAPI{db: _orm}
	post := PostAPI{db: _orm}
	reddit := RedditAPI{db: _orm}
	community := CommunityAPI{db: _orm}

	// Router
	v0 := app.Party("/api/v0", crs).AllowMethods(iris.MethodOptions)

	{
		v0.Get("/", root)

		// auth
		v0.Get("/u/{mail:string}", user.mail)
		v0.Post("/u/{mail:string}", user.verify)

		// open
		// v0.Get("/p/{id:string}", post.spec)
		// v0.Get("/x/{user:string}/p", post.user)

		// author
		v0.Get("/u/{mail:string}/post", post.mail)
		v0.Use(auth)

		// reddit
		v0.Get("/u/:mail/c/:id/reddit", reddit.reddits)
		v0.Post("/u/:mail/reddit", reddit.publish)
		v0.Put("/u/:mail/r/:id", reddit.updateReddit)
		v0.Post("/u/:mail/r/:id/time", reddit.updateRedditTime)
		v0.Delete("/u/:mail/r/:id", reddit.deleteReddit)

		// community
		v0.Get("/u/{mail:string}/c", community.communities)
		v0.Get("/u/{mail:string}/c/:id/members", community.members)
		v0.Get("/u/{mail:string}/c/:id/quit", community.quit)
		v0.Post("/u/{mail:string}/c/create", community.create)
		v0.Post("/u/{mail:string}/c/join", community.join)
		v0.Put("/u/{mail:string}/c/name", community.updateCommunityName)
		// v0.Put("/u/{mail:string}/c/id", community.updateCommunityId)

		// profile
		v0.Post("/u/{mail:string}/upload", user.upload)
		v0.Put("/u/{mail:string}/i/name", user.updateUserName)
		v0.Put("/u/{mail:string}/i/avatar", user.updateUserAvatar)

		// posts
		v0.Get("/u/{mail:string}/post", post.mail)
		v0.Post("/u/{mail:string}/post", user.publish)
		v0.Put("/u/{mail:string}/post/{id:string}", user.updatePost)
		v0.Delete("/u/{mail:string}/post/{id:string}", user.deletePost)
	}

	t := conf()
	app.Run(iris.Addr(t.Get("server.port").(string)))
}

func auth(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	code := ctx.GetHeader("code")

	if rGet(mail) != code {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	ctx.Next()
}
