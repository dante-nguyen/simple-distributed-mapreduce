package server

import "errors"

var (
	errInvalidPort = errors.New("invalid port")
)

type Config struct {
	Port         int
	HostOverride string
}

type OptFunc = func(c *Config)

func (c Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return errInvalidPort
	}

	return nil
}

func defaultConfig(port int) Config {
	return Config{
		Port:         port,
		HostOverride: "",
	}
}

func NewConfig(port int, opts ...OptFunc) (Config, error) {
	res := defaultConfig(port)
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
