package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB() *sqlx.DB {
	dbDriver := "postgres"
	host := "ec2-34-243-101-244.eu-west-1.compute.amazonaws.com"
	port := 5432
	user := "hvbofdxjbkkdgq"
	password := "ff9c8195d4fa5205036cb92a384e142c9ca7bfbbc5f7639f038b4925bacdfea9"
	dbname := "d62omvefcmhpmq"
	dbDataSource := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := sqlx.Open(dbDriver, dbDataSource)
	if err != nil {
		panic(err.Error())
	}
	return db
}
