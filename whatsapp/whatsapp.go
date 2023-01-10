package whatsapp

import (
	"log"
	"os"

	"github.com/macedo/whatsappbot/sqldb"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	"github.com/olebedev/when/rules/br"
	"github.com/olebedev/when/rules/common"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var clients []*Client

var container *sqlstore.Container

var l *log.Logger

var logLevel string

var nlDTParser *when.Parser

type ConnectOptions struct {
	Debug bool
}

func init() {
	logLevel = "INFO"

	nlDTParser = when.New(&rules.Options{
		Distance: 10,
	})
	nlDTParser.Add(common.All...)
	nlDTParser.Add(br.All...)
}

func Connect(opts *ConnectOptions) error {
	if opts.Debug {
		logLevel = "DEBUG"
	}

	l = log.New(os.Stdout, "whatsapp", log.LstdFlags)

	dbLog := waLog.Stdout("DATABASE", logLevel, true)
	container = sqlstore.NewWithDB(sqldb.DB, sqldb.Provider, dbLog)
	if err := container.Upgrade(); err != nil {
		return err
	}

	devices, err := container.GetAllDevices()
	if err != nil {
		return err
	}
	l.Printf("devices found (%d)", len(devices))

	for _, d := range devices {
		if err := ConnectDevice(d); err != nil {
			l.Printf("device-%s could not connect", d.ID)
		}
	}

	return nil
}

func ConnectDevice(d *store.Device) error {
	c := NewClient(d)
	clients = append(clients, c)

	return c.Connect()
}

func Disconnect() {
	for _, c := range clients {
		c.Disconnect()
	}
}

func Clients() []*Client {
	return clients
}
