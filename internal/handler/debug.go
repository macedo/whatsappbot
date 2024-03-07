package handler

import (
	"context"

	"github.com/macedo/whatsappbot/pkg/slice"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type DebugHandler struct {
	cli  *whatsmeow.Client
	jids []string
	log  waLog.Logger
}

type DebugAttributes struct {
	Enabled bool     `mapstructure:"enabled"`
	JIDs    []string `mapstructure:"jids"`
}

func NewDebugHandler(cli *whatsmeow.Client, attrs DebugAttributes) *DebugHandler {
	return &DebugHandler{
		cli:  cli,
		jids: attrs.JIDs,
		log:  waLog.Stdout("debug", "INFO", true),
	}
}

func (h *DebugHandler) EventHandler() func(any) {
	return func(rawEvt any) {
		switch evt := rawEvt.(type) {
		case *events.Message:
			sender := evt.Info.Chat

			if slice.Contains[string](h.jids, sender.String()) {
				_, err := h.cli.SendMessage(context.Background(), sender, &waProto.Message{
					Conversation: proto.String("OK"),
				})

				if err != nil {
					h.log.Errorf("Failed to send message back.\n"+
						"Here's what happened: %v", err)
					return
				}
			}
		}
	}
}
