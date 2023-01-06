package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type TmplData struct {
	H1 string
}

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := &TmplData{
			H1: "Oie",
		}

		err := tmpl.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
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
