package main

import (
	"os"

	"github.com/neuronlabs/neuron/log"
	"github.com/neuronlabs/neuron/repository"

	"github.com/neuronlabs/neuron-extensions/repository/postgres"
)

const (
	envDefaultPostgres  = "NEURON_DEFAULT_POSTGRES"
	envCommentsPostgres = "NEURON_COMMENTS_POSTGRES"
)

// getDefaultPostgresRepository gets the default postgres repository configuration for the service.
func defaultRepository() *postgres.Postgres {
	uriCredentials, ok := os.LookupEnv(envDefaultPostgres)
	if !ok {
		log.Fatalf("required environment variable: %s not found", envDefaultPostgres)
	}
	return postgres.New(repository.WithURI(uriCredentials))
}

// getCommentsPostgresRepository gets the repository configuration for the Comments model.
func commentsRepository() *postgres.Postgres {
	uriCredentials, ok := os.LookupEnv(envCommentsPostgres)
	if !ok {
		log.Fatalf("required environment variable %s not found", envCommentsPostgres)
	}
	return postgres.New(repository.WithURI(uriCredentials))
}
