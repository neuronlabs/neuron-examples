package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/neuronlabs/neuron-extensions/store/memory"

	"github.com/neuronlabs/neuron"
	"github.com/neuronlabs/neuron/auth"
	"github.com/neuronlabs/neuron/log"
	"github.com/neuronlabs/neuron/store"

	"github.com/neuronlabs/neuron-extensions/auth/authenticator"
	"github.com/neuronlabs/neuron-extensions/auth/jwt-tokener"
	serverLogs "github.com/neuronlabs/neuron-extensions/server/http/log"

	"github.com/neuronlabs/neuron-extensions/auth/accounts"
)

func main() {
	// Define the store required for the authenticator.
	inMemoryStore, err := memory.New(store.WithFileName("tmp_store"))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// Define the authenticator that uses SHA256 hash function, with the salt of length 12 bytes
	// and uses in-memory key value store.
	a := authenticator.New(
		auth.AuthenticatorMethod(auth.SHA256),
		auth.AuthenticatorSaltLength(12),
		auth.AuthenticatorStore(inMemoryStore),
	)

	// In order to create authentication tokens auth.Tokener needs to be defined. Let's create the auth.Tokener
	// that takes the secret from the NEURON_TOKENER_SECRET environmental variable, creates the token with hourly expiration time,
	// and refresh tokens with 30 days refresh time.
	secret := os.Getenv("NEURON_TOKENER_SECRET")
	if secret == "" {
		log.Warning("No NEURON_TOKENER_SECRET env found. Setting unsafe secret for the tokener.")
		secret = "secret_generated_in_random"
	}

	tk, err := tokener.New(
		auth.TokenerSigningMethod(jwt.SigningMethodHS256),
		auth.TokenerSecret([]byte(secret)),
		auth.TokenerTokenExpiration(time.Minute*60),
		auth.TokenerRefreshTokenExpiration(time.Minute*60*24*30),
	)
	if err != nil {
		log.Fatalf("Creating tokener failed: %v", err)
	}

	srv, err := getServer()
	if err != nil {
		log.Fatalf("Getting server failed: %v", err)
	}
	log.NewDefault()
	log.SetModulesLevel(log.LevelDebug3)
	n := neuron.New(
		neuron.AccountModel(&accounts.Account{}),
		// Set the authenticator.
		neuron.Authenticator(a),
		// Set the JWT Tokener.
		neuron.Tokener(tk),
		// Set the http server with json:api into service server.
		neuron.Server(srv),
		// Set the default store to be 'in-memory'.
		neuron.DefaultStore(inMemoryStore),
		// Define the default repository name for all models without repository name specified.
		neuron.DefaultRepository(defaultRepository()),
		// Set the models in the service.
		neuron.Models(&accounts.Account{}),
		// Migrate models into service - this would create database definitions for the provided models.
		neuron.MigrateModels(&accounts.Account{}),
		// Initialize model collections that we would like to use
		neuron.Collections(accounts.NRN_Accounts),
	)
	log.SetLevel(log.LevelDebug3)
	serverLogs.SetLevel(log.LevelDebug3)

	ctx := context.Background()
	if err := n.Initialize(ctx); err != nil {
		log.Fatalf("Initialize failed: %v", err)
	}
	defer func() {
		if err := n.CloseAll(ctx); err != nil {
			log.Errorf("Closing failed: %s", err)
		}
	}()

	// List all endpoints defined in the server and print their paths.
	for _, endpoint := range n.Server.GetEndpoints() {
		log.Infof("Endpoint [%s] %s", endpoint.HTTPMethod, endpoint.Path)
	}

	if err := n.Run(ctx); err != nil && err != http.ErrServerClosed {
		log.Errorf("Running neuron service failed: %s", err)
	}
}
