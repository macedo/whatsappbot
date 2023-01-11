package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/macedo/whatsappbot/handlers"
	"github.com/macedo/whatsappbot/scheduler"
	"github.com/macedo/whatsappbot/whatsapp"
)

var debug bool

var l *log.Logger

func init() {
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()
}

func main() {
	l = log.New(os.Stdout, "server", log.LstdFlags)

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.HomePage()).Methods("GET")
	router.HandleFunc("/connect-device", handlers.ConnectDevice(l))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	scheduler.Start()
	defer scheduler.Shutdown()

	opts := &whatsapp.ConnectOptions{Debug: debug}
	if err := whatsapp.Connect(opts); err != nil {
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
