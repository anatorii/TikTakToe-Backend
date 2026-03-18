package main

import (
	"log"
	"tiktaktoe/internal/pkg/di"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found: %v", err)
		return
	}

	(fx.New(di.Module)).Run()
}
