package config

import (
	// "github.com/brumhard/alligotor" I am confused with using this package,
	//  because it is not popular, and low stars in GitHub, so I decided to use other one

	"coupon/internal/app/coupon/api"
	"coupon/internal/pkg/config"
)

type Config struct {
	API     api.Config
	Default config.DefaultConfig
}

func GetConfig() (conf *Config, err error) {
	conf = &Config{}
	err = config.Load(conf)
	return conf, err
}
