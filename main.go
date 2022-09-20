package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xegcrbq/auth/controller"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/repositories"
	"github.com/xegcrbq/auth/services"
	"github.com/xegcrbq/auth/tokenizer"
	"log"
)

func main() {
	//os.Setenv("NODE_ENV", "production")
	//os.Setenv("DB_PORT", "5432")
	//os.Setenv("DB_USER", "postgres")
	//os.Setenv("DB_PASSWORD", "postgres")
	//os.Setenv("DB_DBNAME", "db")
	//os.Setenv("REDIS_HOST", "cache")
	//os.Setenv("REDIS_PORT", "6379")
	//os.Setenv("REDIS_PASSWORD", "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81")
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
