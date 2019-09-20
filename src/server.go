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
	defer _orm.Close()

	user := UserAPI{db: _orm}
	post := PostAPI{db: _orm}
	community := CommunityAPI{db: _orm}

	// Router
	v0 := app.Party("/api/v0", crs).AllowMethods(iris.MethodOptions)

	{
		v0.Get("/", root)

		// auth
		v0.Get("/a/{mail:string}", user.mail)
		v0.Post("/a/{mail:string}", user.verify)

		// open
		v0.Get("/p/{id:string}", post.spec)
		v0.Get("/x/{user:string}/p", post.user)

		v0.Use(auth)

		// community
		v0.Get("/u/{mail:string}/c", community.communities)
		v0.Post("/u/{mail:string}/c/create", community.create)
		v0.Post("/u/{mail:string}/c/join", community.join)

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
