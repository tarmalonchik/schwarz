package main

import (
	"log"
	"runtime"

	"coupon/internal/app/coupon/api"
	"coupon/internal/app/coupon/config"
	configPkg "coupon/internal/pkg/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to load configs: %v", err)
	}

	if cfg.Default.Environment == configPkg.EnvironmentProduction {
		if 32 != runtime.NumCPU() {
			panic("this api is meant to be run on 32 core machines")
		}
	}

	apiServer := api.New(cfg.API)
	log.Println("Starting Coupon service server")
	apiServer.Start()
}
