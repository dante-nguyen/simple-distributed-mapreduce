package worker

import "github.com/caarlos0/env/v11"

type config struct {
	MasterAddr         string `env:"MASTER_ADDR,required,notEmpty"`
	ListenAddr         string `env:"WORKER_LISTEN_ADDR,required,notEmpty"`
	RegisterTimeoutSec uint   `env:"REGISTER_TIMEOUT_SEC,required,notEmpty"`
}

func autoConfig() (config, error) {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}
