package config

import "fmt"

type Environment string

const (
	EnvLocal Environment = "local"
	EnvProd  Environment = "prod"
)

func (e Environment) Validate() error {
	switch e {
	case EnvProd, EnvLocal:
		return nil
	default:
		available := []Environment{EnvLocal, EnvProd}
		return fmt.Errorf("APP_ENV %q is invalid, expected one of %v", e, available)
	}
}
