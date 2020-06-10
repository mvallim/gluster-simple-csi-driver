package config

// Config struct fills the parameters of request or user input
type Config struct {
	Endpoint      string
	NodeID        string
	BlockHostPath string
	Servers       []string
}

//NewConfig returns config struct to initialize new driver
func NewConfig() *Config {
	return &Config{}
}
