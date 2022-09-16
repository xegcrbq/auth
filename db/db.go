package db

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func ConnectDB() *sqlx.DB {
	dbDriver := "postgres"
	db, err := sqlx.Open(dbDriver, os.Getenv("herokuDB1"))
	if err != nil {
		panic(err.Error())
	}
	return db
}
func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:49153",
		Password: "redispw",
		DB:       0,
	})
}
