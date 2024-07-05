package main

import (
	"context"
	"di_container/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Failed to run app: %v", err.Error())
	}
}
