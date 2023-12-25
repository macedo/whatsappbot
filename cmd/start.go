/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/macedo/whatsappbot/internal/handler"
	"github.com/macedo/whatsappbot/internal/storage"
	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// startCmd represents the start command
var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts whatsappbot",
		RunE: func(cmd *cobra.Command, args []string) error {
			log = waLog.Stdout("main", "INFO", true)

			s3Storage := storage.NewS3(newS3Client(), appConfig.Bucket)

			container, err := sqlstore.New("postgres", appConfig.DatabaseURL, waLog.Stdout("database", "INFO", true))
			if err != nil {
				return fmt.Errorf("failed to connect to database %q. \n"+
					"here's what happened': %v", appConfig.DatabaseURL, err)
			}

			device, err := container.GetFirstDevice()
			if err != nil {
				return fmt.Errorf("failed to get device.\n"+
					"here's what happened: %v", err)
			}

			client := whatsmeow.NewClient(device, waLog.Stdout("client", "INFO", true))
			client.PrePairCallback = func(jid types.JID, platform, businessName string) bool {
				log.Infof("pairing %s (platform: %q, business name: %q)", platform, businessName)
				return true
			}

			ch, err := client.GetQRChannel(context.Background())
			if err != nil {
				if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
					log.Errorf("failed to get qr channel from device.\n"+
						"here's what happened: %v", err)
				}
			} else {
				go func() {
					for evt := range ch {
						if evt.Event == "code" {
							qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
						} else {
							log.Infof("qr channel result %s", evt.Event)
						}
					}
				}()
			}

			mediaDownload := handler.NewMediaDownloadHandler(client, appConfig.JIDs, s3Storage, log)
			client.AddEventHandler(mediaDownload.HandlerFunc())

			if err := client.Connect(); err != nil {
				log.Errorf("failed to connect to device %q.\n"+
					"here's what happened: %v", device.ID, err)
				os.Exit(1)
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			<-c
			log.Infof("interrupt received, exiting")
			client.Disconnect()
			log.Infof("bye")

			return nil
		},
	}
)

func newS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Errorf("failed to load aws default config.\n"+
			"here's what happened: %v", err)
		os.Exit(1)
	}

	return s3.NewFromConfig(cfg)
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
