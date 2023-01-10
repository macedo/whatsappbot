package whatsapp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	*whatsmeow.Client
	eventHandlerID uint32
}

func (c *Client) registerHandler() {
	c.eventHandlerID = c.AddEventHandler(c.eventHandler)
}

func (c *Client) eventHandler(evt interface{}) {
	switch evtType := evt.(type) {
	case *events.Message:
		re := regexp.MustCompile(`^(remember me|me lembra de) ?(.*)?$`)
		matches := re.FindStringSubmatch(evtType.Message.GetConversation())
		if len(matches) == 0 {
			l.Println("unhandled message - ", evtType.Message.GetConversation())
			return
		}

		body := matches[0]

		r, err := nlDTParser.Parse(body, time.Now())
		if err != nil {
			l.Fatal(err)
		}

		if r == nil {
			c.SendMessage(context.TODO(), evtType.Info.Sender, &waProto.Message{
				Conversation: proto.String("Desculpe, não entendi direito. Eu sou facilmente confundido. Talvez tente as palavras em uma ordem diferente. Isso geralmente funciona: me lembra de [o que] [quando]"),
			})
		} else {
			reminder := strings.ReplaceAll(body, body[r.Index:r.Index+len(r.Text)], "")
			reminder = strings.TrimSpace(reminder)

			l.Printf("%s - %s", reminder, r.Time)
		}
	}
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
		Client: waCli,
	}
	client.registerHandler()

	return client
}
