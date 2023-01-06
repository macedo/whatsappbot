package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type TmplData struct {
	H1 string
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	data := &TmplData{
		H1: "Oie",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError)
	}
	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			l.Println("read:", err)
			break
		}
		l.Printf("recv: %s", message)
		if err := ws.WriteMessage(mt, message); err != nil {
			l.Printf("write: %s", err)
			break
		}
	}

}

func errorResponse(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, "Ooops", status)
	return
}
