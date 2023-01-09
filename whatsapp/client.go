package whatsapp

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// If you want to access the Client instance inside the event handler, the recommended way is to
// wrap the whole handler in another struct:
//
//	type MyClient struct {
//		WAClient *whatsmeow.Client
//		eventHandlerID uint32
//	}
//
//	func (mycli *MyClient) register() {
//		mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.myEventHandler)
//	}
//
//	func (mycli *MyClient) myEventHandler(evt interface{}) {
//		// Handle event and access mycli.WAClient
//	}
type Client struct {
	*whatsmeow.Client
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

	return &Client{waCli}
}
