package cmd

import (
	"os"

	"github.com/macedo/whatsappbot/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	cfgFile string

	appConfig = &config.AppConfig{}

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
	log = waLog.Stdout("main", "INFO", true)
	store.DeviceProps.Os = proto.String("WhatsApp Bot")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("failed to read config file.\n"+
			"here's what happened: %v", err)
		os.Exit(1)
	}

	log.Infof("Using config file: %s", viper.ConfigFileUsed())

	if err := viper.Unmarshal(appConfig); err != nil {
		log.Errorf("failed to unmarshal config.\n"+
			"here's what happened: %v", err)
		os.Exit(1)
	}
}
