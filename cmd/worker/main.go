package main

import (
	"context"
	"github.com/NikitaVi/microservices_kafka/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app %s", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app %s", err)
	}
}
