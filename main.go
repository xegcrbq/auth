package main

import (
	"github.com/xegcrbq/auth/controller"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/repositories"
	"github.com/xegcrbq/auth/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	sr := repositories.NewSessionRepo(db.ConnectDB())
	cr := repositories.NewCredentialsRepo(db.ConnectDB())
	ss := services.NewSessionService(sr)
	cs := services.NewCredentialsService(cr)
	a := controller.NewAuthController(ss, cs, []byte("djkhgkjdfgndkjnkdjnvkjkdgkjd"))
	app := fiber.New()

	app.Get("/auth/:username-:password", a.Signin)
	app.Get("/", a.Welcome)
	app.Get("/auth/refresh", a.Refresh)

	log.Fatal(app.Listen(":8080"))
}
