package main

import (
	"github.com/xegcrbq/auth/controller"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/repositories"
	"github.com/xegcrbq/auth/services"
	"github.com/xegcrbq/auth/tokenizer"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	sr := repositories.NewSessionRepoRedis(db.ConnectRedis())
	cr := repositories.NewCredentialsRepo(db.ConnectDB())
	ss := services.NewSessionService(sr)
	cs := services.NewCredentialsService(cr)
	service := services.NewService(cs, ss)
	a := controller.NewAuthController(service, tokenizer.NewTestTokenizer())
	app := fiber.New()

	app.Get("/auth/:username-:password", a.Signin)
	app.Get("/", a.Welcome)
	app.Get("/auth/refresh", a.Refresh)

	log.Fatal(app.Listen(":8080"))
}
