package main

import (
	"log"

	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/app"
)

var name, version, compiledAt string //nolint:gochecknoglobals //lld flags

func main() {
	if _, err := app.New(name, version, compiledAt).Start(); err != nil {
		log.Fatal(err)
	}
}
