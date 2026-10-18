package main

import (
	v1 "cryptoapi/internal/endpoint/controller/api/http/v1"
	"cryptoapi/pkg/httpserver"
	"log"
)

func main() {
	chipherController := v1.NewController()

	srv := httpserver.NewServer(chipherController)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
