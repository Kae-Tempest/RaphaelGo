package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func DB() (DB *pgxpool.Pool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//DB, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	DB, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
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
	return DB
}
