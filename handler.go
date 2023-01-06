package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/macedo/whatsappbot/whatsapp"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type HomePageData struct {
	Devices []*store.Device
}

type qrcodeEvent struct {
	Code      string `json:"code"`
	RefreshIn int    `json:"refresh_in"`
	Type      string `json:"type"`
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, nil)
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
		http.Error(w, "Ooops", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	newDevice := whatsapp.NewDevice()
	client := whatsmeow.NewClient(newDevice, waLog.Stdout("DEVICE-NEW", "DEBUG", true))

	qrCh, err := client.GetQRChannel(r.Context())
	if err != nil {
		panic(err)
	}

	if err := client.Connect(); err != nil {
		panic(err)
	}

	for item := range qrCh {
		switch evt := item.Event; evt {
		case "message":
			ws.WriteJSON(&qrcodeEvent{
				Code: item.Code,
				Type: "qrcode",
			})
		default:
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
