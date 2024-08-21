package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	FlagRunAddr         string // address and port to run server
	FlagDBURI           string // driver-specific data source name, usually consisting of at least a database name and connection information.
	FlagAccSysAddr      string // accrual system address
	FlagLogLevel        string
	FlagAttemptInterval int
	FlagRateLimit       int // amount of goroutines being sent to accrual system for recalculation
	PauseDuration       time.Duration
}

func LoadConfig() *Config {
	cfg := new(Config)
	flag.StringVar(&cfg.FlagRunAddr, "a", ":8081", "address and port to run server")
	flag.StringVar(&cfg.FlagDBURI, "d", DefaultPostgresConfig().String(), "DATABASE URI")
	flag.StringVar(&cfg.FlagAccSysAddr, "r", "", "accrual system address")
	flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	flag.IntVar(&cfg.FlagAttemptInterval, "i", 1, "frequency of orders being sent for accrual calculation")
	flag.IntVar(&cfg.FlagRateLimit, "c", 20, "frequency of orders being sent for accrual calculation")
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
	cfg.PauseDuration = time.Duration(cfg.FlagAttemptInterval) * time.Second
	return cfg
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "igortoigildin",
		Password: "Igor109112",
		Database: "postgres",
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
