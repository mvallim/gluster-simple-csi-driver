package config

// Config struct fills the parameters of request or user input
type Config struct {
	Endpoint    string
	NodeID      string
	Servers     string
	HostPath    string
	ServerLabel string
}

//NewConfig returns config struct to initialize new driver
func NewConfig() *Config {
	return &Config{}
}
