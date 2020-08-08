package main

import (
	"context"
	"net/http"

	"github.com/neuronlabs/neuron"
	"github.com/neuronlabs/neuron/log"

	serverLogs "github.com/neuronlabs/neuron-extensions/server/http/log"
)

func main() {
	log.NewDefault()
	n := neuron.New(
		// Set the http server with json:api into service server.
		neuron.Server(getServer()),
		// Define the default repository name for all models without repository name specified.
		neuron.DefaultRepository(defaultRepository()),
		// Set the custom repository for the comments model.
		neuron.RepositoryModels(commentsRepository(), &Comment{}),
		// Set the models in the service.
		neuron.Models(Neuron_Models...),
		// Migrate models into service - this would create database definitions for the provided models.
		neuron.MigrateModels(Neuron_Models...),
		// Initialize model collections that we would like to use
		neuron.Collections(Neuron_Collections...),
	)
	log.SetLevel(log.LevelDebug3)
	serverLogs.SetLevel(log.LevelDebug3)

	ctx := context.Background()
	if err := n.Initialize(ctx); err != nil {
		log.Fatalf("Initialize failed: %v", err)
	}

	// List all endpoints defined in the server and print their paths.
	for _, endpoint := range n.Server.GetEndpoints() {
		log.Infof("Endpoint [%s] %s", endpoint.HTTPMethod, endpoint.Path)
	}
	if err := n.Run(ctx); err != nil && err != http.ErrServerClosed {
		log.Errorf("Running neuron service failed: %s", err)
	}
}
