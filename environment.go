package main

import (
	"os"

	"github.com/rs/zerolog/log"
)

func envString(envName string, defaultValue string) string {
	env, ok := os.LookupEnv(envName)
	if !ok {
		log.Info().Msgf("Environment variable %s is not set, using default value %s", envName, defaultValue)
		return defaultValue
	}
	return env
}
