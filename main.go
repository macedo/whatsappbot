package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var l *log.Logger

var tmpl *template.Template

type TmplData struct {
	H1 string
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	data := &TmplData{
		H1: "Oie",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ooops", http.StatusBadRequest)
		return
	}
}

func main() {
	l = log.New(os.Stdout, "server", log.LstdFlags)

	tmpl = template.Must(template.ParseFiles("index.html"))

	router := mux.NewRouter()
	router.HandleFunc("/", HomePage).Methods("GET")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, syscall.SIGTERM)
	sig := <-sigCh
	l.Println("received terminate, graceful shutdown", sig)

	timeotCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.Shutdown((timeotCtx))
}
