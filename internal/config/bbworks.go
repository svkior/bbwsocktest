package config

type bbConfig struct {
}

func (c *bbConfig) GetBBURL() string {
	return "wss://ropsten1.trezor.io/websocket"
}

func NewBBConfig() *bbConfig {
	return &bbConfig{}
}
