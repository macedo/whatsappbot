package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/macedo/whatsappbot/internal/handler"
	"github.com/macedo/whatsappbot/internal/whatsapp"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts whatsappbot",
		RunE: func(cmd *cobra.Command, args []string) error {
			bot, err := whatsapp.NewBot(appConfig.DB, log)
			if err != nil {
				return err
			}
			defer bot.Disconnect()

			handler.Initialize(appConfig.Handlers, bot)

			if err := bot.Connect(context.Background()); err != nil {
				return err
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			<-c
			log.Infof("interrupt received, exiting")
			log.Infof("bye")

			return nil
		},
	}
)

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
