package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/macedo/whatsappbot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
