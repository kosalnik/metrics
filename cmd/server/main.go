package main

import (
	"metrics/internal/application"
)

func main() {
	app := application.NewApp()
	err := app.Serve()
	if err != nil {
		panic(err)
	}
}
