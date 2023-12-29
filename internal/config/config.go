package config

import "fmt"

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Provider string `mapstructure:"provider"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
}

func (dc *DatabaseConfig) DNS() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		dc.Provider,
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.Name,
	)
}

type HandlerConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	JIDs    []string `mapstructure:"jids"`
}

type AppConfig struct {
	DB       DatabaseConfig `mapstructure:"database"`
	Handlers map[string]any `mapstructure:"handlers"`
}
