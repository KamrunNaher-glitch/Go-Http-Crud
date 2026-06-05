package main

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func connectDb() {
	var err error
	connStr := os.Getenv("DB_STRING")

	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connected Successfully")
}