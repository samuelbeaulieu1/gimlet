package gimlet

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/fatih/structs"
)

const defaultPort string = "8000"
const defaultMaxPostFormSize = 32 << 20

type Config struct {
	Port string

	// Size in MB
	MaxPostFormSizeMB int64

	AllowCredentials bool
	AllowHeaders     []string
	AllowMethods     []string
	AllowOrigin      []string
}

func NewConfig() Config {
	return Config{
		Port:              defaultPort,
		MaxPostFormSizeMB: defaultMaxPostFormSize,
		AllowCredentials:  true,
		AllowHeaders:      []string{},
		AllowMethods:      []string{},
		AllowOrigin:       []string{},
	}
}

func (config *Config) init(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil || os.IsExist(err) {
		content, err := ioutil.ReadFile(filePath)
		if err == nil {
			var parsedConfig map[string]any
			json.Unmarshal(content, &parsedConfig)
			newConfig := structs.Map(config)

			for field, value := range parsedConfig {
				if configValue := newConfig[field]; configValue != nil {
					newConfig[field] = value
				}
			}
			buffer := new(bytes.Buffer)
			json.NewEncoder(buffer).Encode(newConfig)
			json.NewDecoder(buffer).Decode(&config)
			return true
		}
	}
	return false
}
