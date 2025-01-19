package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type DefaultConfig struct {
	CoresNumber uint8       `envconfig:"CORES_NUMBER" default:"32"`
	Environment Environment `envconfig:"ENVIRONMENT" default:"local"`
}

func Load(conf interface{}) error {
	err := envconfig.Process("", conf)
	if err != nil {
		return fmt.Errorf("error env config loading: %w", err)
	}
	return nil
}
