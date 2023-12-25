package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/macedo/whatsappbot/internal/storage"
	"github.com/macedo/whatsappbot/pkg/slice"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type MediaDownloadHandler struct {
	cli     *whatsmeow.Client
	jids    []string
	log     waLog.Logger
	storage storage.Backend
}

func NewMediaDownloadHandler(
	cli *whatsmeow.Client,
	jids string,
	s storage.Backend,
	l waLog.Logger) *MediaDownloadHandler {
	return &MediaDownloadHandler{
		cli:     cli,
		jids:    strings.Split(jids, ","),
		log:     l,
		storage: s,
	}
}

func (h *MediaDownloadHandler) HandlerFunc() func(any) {
	return func(rawEvt any) {
		switch evt := rawEvt.(type) {
		case *events.Message:
			sender := evt.Info.Sender

			if slice.Contains[string](h.jids, sender.String()) {
				document := evt.Message.GetDocumentMessage()
				if document != nil {
					filename := strings.ReplaceAll(*document.FileName, " ", "_")
					data, err := h.cli.Download(document)
					if err != nil {
						h.log.Errorf("Failed to download document %q. \n"+
							"Here's what happened: %v", *document.FileName, err)
						return
					}

					path := fmt.Sprintf("%s/%s", sender.String(), filename)

					if err := h.storage.Save(context.Background(), path, data); err != nil {
						h.log.Errorf("Failed to save document %q.\n"+
							"Here's what happened: %v", *document.FileName, err)
					}

					h.log.Infof("uploaded %q", path)
				}
			}
		}
	}
}
