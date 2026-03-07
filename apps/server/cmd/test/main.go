package main

import (
	"context"
	"log"
	"os"

	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	blob "github.com/aidan-neel/shulker/apps/server/src/blob"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)
	queries := db.New(sqlDB)

	ctx := context.Background()
	blobRepo := blob.NewBlobRepo(queries)
	blobService := blob.NewService(blobRepo)

	data, err := os.ReadFile("./5MB.json")
	if err != nil {
		panic(err)
	}

	blob, err := blobService.Put(ctx, data, "019cc122-be83-74bd-a725-6aae5019bbd1", "/images/hey.png", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	data, err = blobService.Get(ctx, blob.Hash)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("./test.json", data, 0644)
}
