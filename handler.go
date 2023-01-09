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

type event struct {
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

	client := whatsapp.NewClient(nil)

	qrCh, err := client.GetQRChannel(r.Context())
	if err != nil {
		l.Println("error requesting for new qrcode - ", err)
		ws.Close()
	}

	if err := client.Connect(); err != nil {
		l.Println("error connecting new device - ", err)
		ws.Close()
	}

	for item := range qrCh {
		l.Printf("event %s", item.Event)
		switch evt := item.Event; evt {
		case "code":
			ws.WriteJSON(&qrcodeEvent{
				Code: item.Code,
				Type: "qrcode",
			})

		case "success":
			l.Printf("connected new device - %s", client.Store.ID)
			whatsapp.ConnectDevice(client.Store)
			client.Disconnect()
			ws.WriteJSON(&event{
				Type: "success",
			})
			ws.Close()

		default:
			ws.WriteJSON(&event{
				Type: evt,
			})
			ws.Close()
		}
	}
}

func renderPage(w io.Writer, data any) error {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}
