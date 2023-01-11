package whatsapp

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Client struct {
	*whatsmeow.Client
	eventHandlerIDs []uint32
}

func (c *Client) registerHandler() {
	c.eventHandlerIDs = append(c.eventHandlerIDs, c.AddEventHandler(MessagesHandler(c)))
}

func NewClient(d *store.Device) *Client {
	var id string
	if d == nil {
		d = container.NewDevice()
		id = "NEW"
	} else {
		id = d.ID.String()
	}

	cliLog := waLog.Stdout(fmt.Sprintf("DEVICE-%s", id), logLevel, true)
	waCli := whatsmeow.NewClient(d, cliLog)

	client := &Client{
		Client:          waCli,
		eventHandlerIDs: []uint32{},
	}
	client.registerHandler()

	return client
}
