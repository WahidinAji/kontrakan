package main

import (
	"os"

	"github.com/rs/zerolog/log"
)

func envString(envName, defaultValue string) (env string) {
	env, ok := os.LookupEnv(envName)
	if !ok {
		log.Info().Msgf("Environment variable %s is not set, using default value %s", envName, defaultValue)
		return defaultValue
	}
	return
}

func setEnvDB(envName, defaultEnv string) (db string) {
	db, ok := os.LookupEnv(envName)
	if !ok {
		log.Info().Msgf("Environment variable is not set, using default value %s", defaultEnv)
		db = defaultEnv
		return
	}

	return
}
