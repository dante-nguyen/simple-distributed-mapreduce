package master

import "github.com/caarlos0/env/v11"

type config struct {
	ListenAddress string `env:"MASTER_LISTEN_ADDR,required,notEmpty"`
}

func autoConfig() (config, error) {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}
