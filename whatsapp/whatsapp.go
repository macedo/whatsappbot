package whatsapp

import (
	"fmt"

	"github.com/macedo/whatsappbot/sqldb"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var clients map[string]*whatsmeow.Client

var container *sqlstore.Container

var Devices []*sqlstore.Container

func Connect() error {
	container = sqlstore.NewWithDB(sqldb.DB, "sqlite3", nil)

	devices, err := container.GetAllDevices()
	if err != nil {
		return err
	}

	for _, device := range devices {
		id := device.ID.String()

		clients[id] = whatsmeow.NewClient(
			device,
			waLog.Stdout(fmt.Sprintf("DEVICE-%s", id), "DEBUG", true),
		)
		go func() {
			clients[id].AddEventHandler(echoHandler)
			err := clients[id].Connect()
			if err != nil {
				panic(err)
			}
		}()
	}

	return nil
}

func NewDevice() *store.Device {
	return container.NewDevice()
}

func echoHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received Message", v.Message.GetConversation())
	}
}
