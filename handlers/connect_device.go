package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/macedo/whatsappbot/whatsapp"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type event struct {
	Type string `json:"type"`
}

type qrcodeEvent struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

func ConnectDevice(l *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.Printf("error during connection upgrade - %s", err)
			return
		}
		defer ws.Close()

		client := whatsapp.NewClient(nil)

		qrCh, err := client.GetQRChannel(r.Context())
		if err != nil {
			l.Printf("error requesting for new qrcode - %s", err)
			ws.Close()
		}

		if err := client.Connect(); err != nil {
			l.Printf("error connecting new device - %s", err)
			ws.Close()
		}

		for item := range qrCh {
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
					Type: evt,
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
}
