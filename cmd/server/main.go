package main

import (
	"github.com/kosalnik/metrics/internal/server"
)

func main() {
	app := server.NewApp()
	err := app.Serve()
	if err != nil {
		panic(err)
	}
}
