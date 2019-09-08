package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Logger() //.SetLevel("debug")
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
	article := ArticleAPI{db: _orm}

	// Router
	v0 := app.Party("/api/v0", crs).AllowMethods(iris.MethodOptions)

	{
		v0.Get("/", root)
		v0.Get("/{mail:string}/code", user.sendCode)
		v0.Post("/{mail:string}/verify", user.verify)
		v0.Get("/{mail:string}/articles", article.mail)
		v0.Get("/articles/{user:string}", article.user)
		v0.Get("/article/{id:string}", article.spec)

		v0.Use(auth)

		v0.Post("/{mail:string}/publish", user.publish)
		v0.Post("/{mail:string}/upload", user.upload)
		v0.Post("/{mail:string}/update/name", user.updateUserName)
		v0.Post("/{mail:string}/article/update", user.updateArticle)
		v0.Post("/{mail:string}/article/delete", user.deleteArticle)

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
