// Package config
// nolint
//
//go:generate go-enum -f=$GOFILE --nocase
package config

// Environment ENUM(local, dev, staging, production)
type Environment string
