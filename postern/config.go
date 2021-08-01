package postern

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	LocalPort  int		`json:"local_port"`
	RemotePort int		`json:"remote_port"`
	RemoteIp   string	`json:"remote_ip"`
	Key        string	`json:"key"`
	Method	   int		`json:"cipher"`
}

func LoadConfig(path string) (config *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	jsonStr, err := io.ReadAll(file)
	if err != nil {
		return
	}

	config = &Config{}
	if err = json.Unmarshal(jsonStr, config); err != nil {
		return nil, err
	}
	return
}
