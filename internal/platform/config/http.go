package config

import (
	"errors"
	"net"
	"strconv"
	"time"
)

type HTTPConfig struct {
	Host              string        `env:"HOST" envDefault:"0.0.0.0"`
	Port              uint16        `env:"PORT" envDefault:"8080"`
	ReadyTimeout      time.Duration `env:"READY_TIMEOUT" envDefault:"2s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
}

func (hc HTTPConfig) Addr() string {
	portStr := strconv.FormatUint(uint64(hc.Port), 10)

	return net.JoinHostPort(hc.Host, portStr)
}

func (hc HTTPConfig) validate() error {
	if hc.Port == 0 {
		return errors.New("HTTP_PORT can not be 0")
	}

	if hc.ReadyTimeout <= 0 {
		return errors.New("HTTP_READY_TIMEOUT must be greater than 0")
	}

	if hc.ReadTimeout <= 0 {
		return errors.New("HTTP_READ_TIMEOUT must be greater than 0")
	}

	if hc.WriteTimeout <= 0 {
		return errors.New("HTTP_WRITE_TIMEOUT must be greater than 0")
	}

	if hc.ReadHeaderTimeout <= 0 {
		return errors.New("HTTP_READ_HEADER_TIMEOUT must be greater than 0")
	}

	if hc.IdleTimeout <= 0 {
		return errors.New("HTTP_IDLE_TIMEOUT must be greater than 0")
	}

	return nil
}
