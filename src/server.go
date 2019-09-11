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

	// Router
	v0 := app.Party("/api/v0", crs).AllowMethods(iris.MethodOptions)

	{
		v0.Get("/", root)

		// auth
		v0.Get("/a/{mail:string}", user.mail)    // ok
		v0.Post("/a/{mail:string}", user.verify) // ok

		// open
		v0.Get("/p/{id:string}", post.spec)
		v0.Get("/x/{user:string}/p", post.user)

		v0.Use(auth) // ok

		// profile
		v0.Post("/u/{mail:string}/upload", user.upload)        // ok
		v0.Put("/u/{mail:string}/i/name", user.updateUserName) // ok

		// posts
		v0.Get("/u/{mail:string}/p", post.mail)                      // ok
		v0.Post("/u/{mail:string}/p", user.publish)                  // ok
		v0.Put("/u/{mail:string}/p/{id:string}", user.updatePost)    // ok
		v0.Delete("/u/{mail:string}/p/{id:string}", user.deletePost) // ok
	}

	app.Run(iris.Addr(":6060"))
}

func auth(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	code := ctx.GetHeader("code")

	if rGet(mail) != code {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	ctx.Next()
}
