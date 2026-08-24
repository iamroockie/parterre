package config

import "fmt"

type Env string

const (
	EnvProd  Env = "prod"
	EnvLocal Env = "local"
)

func (e Env) String() string {
	return string(e)
}

func (e *Env) UnmarshalText(data []byte) error {
	format := Env(string(data))

	switch format {
	case EnvProd, EnvLocal:
		*e = format
		return nil
	default:
		return fmt.Errorf("unknown env %q", format)
	}
}
