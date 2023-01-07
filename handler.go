package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/macedo/whatsappbot/whatsapp"
)

type homePageData struct {
	Clients []*whatsapp.Client
}

type qrcodeEvent struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

type timeoutEvent struct {
	Type string `json:"type"`
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	data := &homePageData{
		Clients: whatsapp.Clients(),
	}
	err := renderPage(w, data)
	if err != nil {
		http.Error(w, "Ooops", http.StatusBadRequest)
		return
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
		l.Println("error during connection upgrade - ", err)
		return
	}
	defer ws.Close()

	// client := whatsapp.NewClient()

	// qrCh, err := client.GetQRChannel(r.Context())
	// if err != nil {
	// 	panic(err)
	// }

	// client.Log = waLog.Stdout("asd", "DEBUG", true)

	// if err := client.Connect(); err != nil {
	// 	panic(err)
	// }

	// for item := range qrCh {
	// 	l.Printf("event %s", item.Event)
	// 	switch evt := item.Event; evt {
	// 	case "code":
	// 		ws.WriteJSON(&qrcodeEvent{
	// 			Code: item.Code,
	// 			Type: "qrcode",
	// 		})
	// 	case "success":

	// 	case "timeout":
	// 		ws.WriteJSON(&timeoutEvent{
	// 			Type: "timeout",
	// 		})
	// 		ws.Close()
	// 	default:
	// 	}
	// }

}

func renderPage(w io.Writer, data any) error {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}
