package main

import (
	"log"
	"os"

	_ "github.com/vedomirr/remindista/docs"
	"github.com/vedomirr/remindista/internal/app"
)

// @title remindista
// @version 0.0.1
// @description remindista API

// @host 10.0.0.109:8888
// @basePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(a.Run())
}
