package whatsapp

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/macedo/whatsappbot/internal/config"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Bot struct {
	Client *whatsmeow.Client
	device *store.Device
	log    waLog.Logger
}

func NewBot(dbConfig config.DatabaseConfig, log waLog.Logger) (*Bot, error) {
	container, err := sqlstore.New(dbConfig.Provider, dbConfig.DNS(), waLog.Stdout("database", "INFO", true))
	if err != nil {
		return nil, fmt.Errorf("failed to connecto to database %q. \n"+
			"here's what happened: %v", dbConfig.DNS(), err)
	}

	device, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device.\n"+
			"here's what happened: %v", err)
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("client", "INFO", true))

	return &Bot{
		Client: client,
		device: device,
		log:    log,
	}, nil

}

func (bot *Bot) Connect(ctx context.Context) error {
	ch, err := bot.Client.GetQRChannel(ctx)
	if err != nil {
		if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			bot.log.Errorf("failed to get QR channel from device.\n"+
				"here's what happened: %v", err)
		}
	} else {
		go func() {
			for item := range ch {
				if item.Event == "code" {
					qrterminal.GenerateHalfBlock(item.Code, qrterminal.L, os.Stdout)
				} else {
					bot.log.Infof("QR channel result %s,", item.Event)
				}
			}
		}()
	}

	if err := bot.Client.Connect(); err != nil {
		return fmt.Errorf("failed to connect to device %q.\n"+
			"here's what happened: %v", bot.device.ID, err)
	}

	return nil
}

func (bot *Bot) Disconnect() {
	bot.Client.Disconnect()
}
