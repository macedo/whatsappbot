package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var (
	cfgFile string

	appConfig = &AppConfig{}

	log waLog.Logger

	rootCmd = &cobra.Command{
		Use:   "whatsappbot",
		Short: "WhatsApp Bot",
		Long:  "WhatsApp Bot is a CLI that expand WhatsApp capabilities with custom message handlers",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is app.env)")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("app")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("failed to read config file.\n"+
			"here's what happened: %v", err)
		os.Exit(1)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if err := viper.Unmarshal(appConfig); err != nil {
		log.Errorf("failed to unmarshal config.\n"+
			"here's what happened: %v", err)
		os.Exit(1)
	}
}
