package config

type Config struct {
	Endpoint string
	NodeID   string
}

func NewConfig() *Config {
	return &Config{}
}
