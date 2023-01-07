package whatsapp

import (
	"log"
	"os"

	"github.com/macedo/whatsappbot/sqldb"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

var clients []*Client

var container *sqlstore.Container

var l *log.Logger

func init() {
	l = log.New(os.Stdout, "whatsapp", log.LstdFlags)

	container = sqlstore.NewWithDB(sqldb.DB, "sqlite3", nil)
	if err := container.Upgrade(); err != nil {
		log.Fatal(err)
	}
}

func Connect() error {
	devices, err := container.GetAllDevices()
	if err != nil {
		return err
	}
	l.Printf("devices found (%d)", len(devices))

	for _, d := range devices {
		c := NewClient(d)
		clients = append(clients, c)

		if err := c.Connect(); err != nil {
			l.Printf("device-%s could not connect", c.Store.ID.String())
		}
	}

	return nil
}

func Disconnect() {
	for _, c := range clients {
		c.Disconnect()
	}
}

func Clients() []*Client {
	return clients
}

func Devices() []*store.Device {
	devices := []*store.Device{}

	for _, cli := range clients {
		devices = append(devices, cli.Store)
	}

	return devices
}
