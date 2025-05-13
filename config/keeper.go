package config

import (
	"bytes"
	"os"

	"github.com/humangrass/gommon/database"
	"github.com/humangrass/price-keeper/pgk/x/xhttp"
	"gopkg.in/yaml.v3"
)

type Keeper struct {
	Workers      uint         `yaml:"workers"`
	Runners      any          `yaml:"runners"`
	IsProduction bool         `yaml:"is_production"`
	Server       xhttp.Opt    `yaml:"server"`
	Database     database.Opt `yaml:"database"`
}

func NewKeeperConfig(filepath string) (*Keeper, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Keeper
	data := yaml.NewDecoder(bytes.NewReader(content))
	err = data.Decode(&cfg)

	return &cfg, err
}
