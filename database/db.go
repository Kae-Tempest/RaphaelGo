package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func DB() (DB *pgx.Conn) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to connect to databases: %v\n", err)
		if err != nil {
			_ = fmt.Errorf("error while send error msh for connection database: %s", err)
			return nil
		}
		os.Exit(1)
	} else {
		fmt.Println("Database Connected !")
	}
	defer DB.Close(context.Background())

	return DB
}
