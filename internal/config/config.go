package config

import (
	"errors"
	"flag"
	"log"
	"os"
)

var (
	ErrNoRunAddress = errors.New("run address not set in env and cla")
	ErrNoDBURI      = errors.New("db uri not set in env and cla")
	ErrAccrualSA    = errors.New("accrual system address not set in env and cla")
)

// Config -.
type Config struct {
	RunAddress           string // Адрес и порт запуска сервиса
	DataBaseUri          string // Адрес подключения к базе данных:
	AccrualSystemAddress string // Адрес системы расчёта начислений
}

// New returns app config
func New() (*Config, error) {
	cfg := &Config{}

	cfg.RunAddress = os.Getenv("RUN_ADDRESS")
	if cfg.RunAddress == "" {
		log.Println("config info: server address not set in env")
		RunAddressFromCMA(cfg)
		if cfg.RunAddress == "" {
			log.Println("config info: server address not set in cla")
			return cfg, ErrNoRunAddress
		}
	}

	cfg.DataBaseUri = os.Getenv("DATABASE_URI")
	if cfg.DataBaseUri == "" {
		log.Println("config info: db uri not set in env")
		DataBaseURIFromCMA(cfg)
		if cfg.DataBaseUri == "" {
			log.Println("config info: db uri not set in cla")
			return cfg, ErrNoDBURI
		}
	}

	cfg.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if cfg.AccrualSystemAddress == "" {
		log.Println("config info: accrual system address not set in env")
		AccrualSystemAddressFromCMA(cfg)
		if cfg.AccrualSystemAddress == "" {
			log.Println("config info: accrual system address not set in cla")
			return cfg, ErrAccrualSA
		}
	}

	flag.Parse()

	return cfg, nil
}

// RunAddressFromCMA -.
func RunAddressFromCMA(cfg *Config) {
	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "host(server address) to listen on")
}

// DataBaseURIFromCMA -.
func DataBaseURIFromCMA(cfg *Config) {
	flag.StringVar(&cfg.DataBaseUri, "d", "postgres://postgres:postgres@172.17.0.2:5432/postgres?sslmode=disable", "db destination")
}

// AccrualSystemAddressFromCMA -.
func AccrualSystemAddressFromCMA(cfg *Config) {
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "localhost:8080", "accrual system address")
}
