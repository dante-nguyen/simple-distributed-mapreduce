package server

import "errors"

var (
	errInvalidPort        = errors.New("invalid port")
	errEmptyAdvertiseAddr = errors.New("empty advertise address")
)

type Config struct {
	Port          int
	HostOverride  string
	AdvertiseAddr string
}

type OptFunc = func(c *Config)

func (c Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return errInvalidPort
	}

	if len(c.AdvertiseAddr) == 0 {
		return errEmptyAdvertiseAddr
	}

	return nil
}

func defaultConfig(port int, advertiseAddr string) Config {
	return Config{
		Port:          port,
		HostOverride:  "",
		AdvertiseAddr: advertiseAddr,
	}
}

func NewConfig(port int, advertiseAddr string, opts ...OptFunc) (Config, error) {
	res := defaultConfig(port, advertiseAddr)
	for _, opt := range opts {
		opt(&res)
	}

	if err := res.validate(); err != nil {
		return Config{}, err
	}

	return res, nil
}

// options

func WithHostOverride(hostname string) OptFunc {
	return func(c *Config) {
		c.HostOverride = hostname
	}
}
