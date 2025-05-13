package main

import (
	"context"
	"log"
	"simple-finance/internal/app"
)

func main() {
	ctx := context.Background()

	app, err := app.NewApp(ctx)

	if err != nil {
		log.Fatal(err)

	}
	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
