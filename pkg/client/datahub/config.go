package datahub

import (
	"time"
)

// Config provides a configuration struct to connect to Alameda-Datahub
type Config struct {
	Timeout time.Duration
	Address string `mapstructure:"address"`
}

// NewDefaultConfig returns default configuration
func NewDefaultConfig() Config {
	return Config{
		Timeout: 30 * time.Second,
		Address: "alameda-datahub.federatorai.svc.cluster.local:50050",
	}
}
