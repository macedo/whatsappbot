package cmd

type AppConfig struct {
	Bucket      string `mapstructure:"BUCKET"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JIDs        string `mapstructure:"JIDS"`
}
