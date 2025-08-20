package main

import (
	"context"
	"log"

	"github.com/G0SU19O2/Go-Twitter/config"
	"github.com/G0SU19O2/Go-Twitter/postgres"
)

func main() {
	ctx := context.Background()
	conf := config.New()
	db := postgres.New(ctx, conf)
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}
}
