package worker

import (
	"errors"
	"time"
)

var (
	errEmptyMasterAddr    = errors.New("empty master address")
	errEmptyAdvertiseAddr = errors.New("empty advertise address")
)

type Config struct {
	MasterAddr      string
	AdvertiseAddr   string
	RegisterTimeout time.Duration
}

func validateConfig(cfg Config) error {
	if len(cfg.MasterAddr) == 0 {
		return errEmptyMasterAddr
	} else if len(cfg.AdvertiseAddr) == 0 {
		return errEmptyAdvertiseAddr
	}

	return nil
}
