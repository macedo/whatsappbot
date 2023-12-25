package main

import (
	"log"

	"github.com/macedo/whatsappbot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
