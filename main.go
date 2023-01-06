package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Oie"))
	})

	addr := ":8080"

	fmt.Printf("listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed")
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
