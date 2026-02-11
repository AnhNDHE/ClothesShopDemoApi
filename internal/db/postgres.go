package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Connect initializes the PostgreSQL connection
func Connect(dbUrl string) {
	var err error
	DB, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("❌ Connect DB error:", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("❌ Ping DB error:", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
