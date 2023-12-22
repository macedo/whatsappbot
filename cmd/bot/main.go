package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"
	"github.com/macedo/whatsappbot/internal/handler"
	"github.com/macedo/whatsappbot/internal/storage"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	cfg, err := LoadDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	jids := strings.Split(cfg.JIDs, ",")

	log := waLog.Stdout("main", "INFO", true)

	s3Storage := storage.NewS3(newS3Client(), cfg.Bucket)

	container, err := sqlstore.New("postgres", cfg.DatabaseURL, waLog.Stdout("database", "INFO", true))
	if err != nil {
		log.Errorf("Failed to connect to database %q. \n"+
			"Here's what happened: %v", cfg.DatabaseURL, err)
		os.Exit(1)
	}

	device, err := container.GetFirstDevice()
	if err != nil {
		log.Errorf("Failed to get device: %q. \n"+
			"Here's what happened: %v", err)
		os.Exit(1)
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("client", "INFO", true))
	client.PrePairCallback = func(jid types.JID, platform, businessName string) bool {
		log.Infof("Pairing %s (platform: %q, business name: %q)", platform, businessName)
		return true
	}

	ch, err := client.GetQRChannel(context.Background())
	if err != nil {
		if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			log.Errorf("Failed to get QR channel from device.\n"+
				"Here's what happened: %v", err)
		}
	} else {
		go func() {
			for evt := range ch {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else {
					log.Infof("QR channel result: %s", evt.Event)
				}
			}
		}()
	}

	mediaDownload := handler.NewMediaDownloadHandler(client, jids, s3Storage, log)
	client.AddEventHandler(mediaDownload.HandlerFunc())

	if err := client.Connect(); err != nil {
		log.Errorf("Failed to connect to device. %q \n"+
			"Here's what happened: %v", device.ID, err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Infof("Interrupt received, exiting")
	client.Disconnect()
	log.Infof("Bye")
}

func newS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
}
