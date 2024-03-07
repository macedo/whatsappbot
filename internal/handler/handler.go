package handler

import (
	"log"

	"github.com/macedo/whatsappbot/internal/whatsapp"
	"github.com/mitchellh/mapstructure"
)

func Initialize(handlers map[string]any, bot *whatsapp.Bot) {
	for k, v := range handlers {
		switch k {
		case "debug":
			var attrs DebugAttributes
			if err := mapstructure.Decode(v, &attrs); err != nil {
				log.Fatal(err)
			}
			h := NewDebugHandler(bot.Client, attrs)

			if attrs.Enabled {
				bot.Client.AddEventHandler(h.EventHandler())
			}

		case "media_download":
			var attrs MediaDownloadAttributes
			if err := mapstructure.Decode(v, &attrs); err != nil {
				log.Fatal(err)
			}
			h := NewMediaDownloadHandler(bot.Client, attrs)

			if attrs.Enabled {
				bot.Client.AddEventHandler(h.EventHandler())
			}
		}
	}
}
