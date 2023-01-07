package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/macedo/whatsappbot/whatsapp"
)

var l *log.Logger

func main() {
	l = log.New(os.Stdout, "server", log.LstdFlags)

	router := mux.NewRouter()
	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/ws", WS)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	if err := whatsapp.Connect(); err != nil {
		l.Fatal(err)
	}

	go func() {
		l.Printf("listening on %s", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	l.Println("received terminate, graceful shutdown", sig)

	timeotCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	whatsapp.Disconnect()
	srv.Shutdown((timeotCtx))
}
