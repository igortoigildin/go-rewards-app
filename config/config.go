package config

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr    string
	FlagDBURI      string
	FlagAccSysAddr string
	FlagLogLevel   string
}

func LoadConfig() *Config {
	cfg := new(Config)
	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.FlagDBURI, "d", "", "DATABASE URI")
	flag.StringVar(&cfg.FlagAccSysAddr, "r", "", "accrual system address")
	flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	flag.Parse()
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		cfg.FlagRunAddr = envRunAddr
	}
	if envDBURI := os.Getenv("DATABASE_URI"); envDBURI != "" {
		cfg.FlagDBURI = envDBURI
	}
	if envAccSysAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccSysAddr != "" {
		cfg.FlagAccSysAddr = envAccSysAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}
	return cfg
}
