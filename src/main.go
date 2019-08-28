package main

import (
	"github.com/kataras/iris"

	"github.com/iris-contrib/middleware/cors"
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
	user := User{}
	article := Article{}

	// Router
	v0 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)

	{
		v0.Get("/", root)
		v0.Post("/{mail:string}/code", user.sendCode)
		v0.Post("/{mail:string}/verify", user.verify)
		v0.Post("/{mail:string}/publish", user.publish)
		v0.Get("/{mail:string}/articles", article.articles)
	}

	app.Run(iris.Addr(":8080"))
}