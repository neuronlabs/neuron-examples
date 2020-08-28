package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/neuronlabs/neuron-extensions/codec/json"

	stdHttp "net/http"

	"github.com/neuronlabs/neuron-extensions/server/http"
	"github.com/neuronlabs/neuron-extensions/server/http/api/authentication"
	"github.com/neuronlabs/neuron-extensions/server/http/httputil"
	"github.com/neuronlabs/neuron-extensions/server/http/middleware"
	"github.com/neuronlabs/neuron/auth"
	"github.com/neuronlabs/neuron/codec"
	"github.com/neuronlabs/neuron/core"
	"github.com/neuronlabs/neuron/log"
	"github.com/neuronlabs/neuron/mapping"
	"github.com/neuronlabs/neuron/server"

	"github.com/neuronlabs/neuron-extensions/auth/accounts"
)

type tokenChecker struct {
	c *core.Controller
}

func (l *tokenChecker) GetEndpoints() []*server.Endpoint {
	return []*server.Endpoint{{
		Path: "/auth/verify-token",
	}}
}

func (l *tokenChecker) InitializeAPI(c *core.Controller) error {
	l.c = c
	return nil
}

func (l *tokenChecker) SetRoutes(router *httprouter.Router) error {
	chain := server.MiddlewareChain{middleware.Controller(l.c), middleware.WithCodec(json.GetCodec(l.c)), middleware.BearerAuthenticate()}
	router.GET("/auth/verify-token", httputil.Wrap(chain.Handle(stdHttp.HandlerFunc(func(rw stdHttp.ResponseWriter, req *stdHttp.Request) {
		acc, ok := auth.CtxGetAccount(req.Context())
		if !ok {
			rw.WriteHeader(500)
			rw.Write([]byte("Account not found"))
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(stdHttp.StatusOK)
		if err := json.GetCodec(l.c).MarshalPayload(rw, &codec.Payload{
			ModelStruct: l.c.MustModelStruct(acc),
			Data:        []mapping.Model{acc},
		}, codec.MarshalSingleModel()); err != nil {
			log.Errorf("Marshaling payload failed: %v", err)
		}
	}))))
	return nil
}

var _ http.API = &tokenChecker{}

func getAuthAPI() (http.API, error) {
	api, err := authentication.New(
		authentication.WithAccountModel(&accounts.Account{}),
		authentication.WithPathPrefix("/auth"),
	)
	return api, err
}

func getServer() (*http.Server, error) {
	// Create api based on json:api specification.
	authAPI, err := getAuthAPI()
	if err != nil {
		return nil, err
	}
	// Create new http server.
	s := http.New(
		http.WithAPI(authAPI),
		http.WithAPI(&tokenChecker{}),
		// Mount json:api with the models.
		// Set the listening port.
		http.WithPort(8080),
	)
	return s, nil
}
