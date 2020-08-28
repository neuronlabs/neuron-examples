package main

import (
	"github.com/neuronlabs/neuron-extensions/server/xhttp"
	"github.com/neuronlabs/neuron-extensions/server/xhttp/api/jsonapi"
	"github.com/neuronlabs/neuron-extensions/server/xhttp/middleware"
)

func getJsonAPI() *jsonapi.API {
	// Define new json:api specification API.
	api := jsonapi.New(
		jsonapi.WithMiddlewares(middleware.ResponseWriter(-1), middleware.LogRequest),
		// Set path prefix to '/v1/api/
		jsonapi.WithPathPrefix("/v1/api"),
		// Set the strict unmarshal flag which disallows unknown model fields in the documents.
		jsonapi.WithStrictUnmarshal(),
		// Set the models with default api handler.
		jsonapi.WithDefaultHandlerModels(&Post{}, &Comment{}),
		// Set the Blog model with the blog handler.
		jsonapi.WithModelHandler(&Blog{}, BlogHandler{}),
	)
	return api
}

func getServer() *http.Server {
	// Create api based on json:api specification.
	jsonAPI := getJsonAPI()
	// Create new http server.
	s := http.New(
		// Mount json:api with the models.
		http.WithAPI(jsonAPI),
		// Set the listening port.
		http.WithPort(8080),
	)
	return s
}
