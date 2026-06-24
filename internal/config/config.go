package config

import "os"

type Config struct {
	Gpp     string
	Clangpp string
}

func Load() Config {
	return Config{
		Gpp:     os.Getenv("CPX_GPP"),
		Clangpp: os.Getenv("CPX_CLANGPP"),
	}
}
