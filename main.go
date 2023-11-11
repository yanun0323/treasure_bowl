package main

import (
	"log"
	"main/internal/app"
	"main/pkg/infra"
)

func main() {
	if err := infra.Init("config"); err != nil {
		log.Fatalf("init infra, err: %+v", err)
	}
	app.Run()
}
