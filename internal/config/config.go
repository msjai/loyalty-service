package config

import (
	"errors"
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/msjai/loyalty-service/internal/logger"
)

var (
	ErrNoRunAddress = errors.New("config - New - RunAddressFromCMA: run address not set in env and cla")
	ErrNoDBURI      = errors.New("config - New - DataBaseURIFromCMA: db uri not set in env and cla")
	ErrAccrualSA    = errors.New("config - New - AccrualSystemAddressFromCMA: accrual system address not set in env and cla")
)

// Config -.
type Config struct {
	RunAddress           string             // Адрес и порт запуска сервиса
	DataBaseURI          string             // Адрес подключения к базе данных:
	AccrualSystemAddress string             // Адрес системы расчёта начислений
	L                    *zap.SugaredLogger // Логгер
}

// New returns app config
func New() (*Config, error) {
	cfg := &Config{}

	cfg.L = logger.New()

	cfg.RunAddress = os.Getenv("RUN_ADDRESS")
	if cfg.RunAddress == "" {
		cfg.L.Infow("config info: server address not set in env")
		RunAddressFromCMA(cfg)
		if cfg.RunAddress == "" {
			cfg.L.Infow("config info: server address not set in cla")
			return cfg, ErrNoRunAddress
		}
	}

	cfg.DataBaseURI = os.Getenv("DATABASE_URI")
	if cfg.DataBaseURI == "" {
		cfg.L.Infow("config info: db uri not set in env")

		DataBaseURIFromCMA(cfg)
		if cfg.DataBaseURI == "" {
			cfg.L.Infow("config info: db uri not set in cla")
			return cfg, ErrNoDBURI
		}
	}

	cfg.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if cfg.AccrualSystemAddress == "" {
		cfg.L.Infow("config info: accrual system address not set in env")
		AccrualSystemAddressFromCMA(cfg)
		if cfg.AccrualSystemAddress == "" {
			cfg.L.Infow("config info: accrual system address not set in cla")
			return cfg, ErrAccrualSA
		}
	}

	flag.Parse()

	return cfg, nil
}

// RunAddressFromCMA -.
func RunAddressFromCMA(cfg *Config) {
	flag.StringVar(&cfg.RunAddress, "a", "", "host(server address) to listen on")
}

// DataBaseURIFromCMA -.
func DataBaseURIFromCMA(cfg *Config) {
	flag.StringVar(&cfg.DataBaseURI, "d", "", "db destination")
}

// AccrualSystemAddressFromCMA -.
func AccrualSystemAddressFromCMA(cfg *Config) {
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "accrual system address")
}
