module github.com/neuronlabs/neuron-examples/server/authentication

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/julienschmidt/httprouter v1.3.0
	github.com/neuronlabs/neuron v0.19.1
	github.com/neuronlabs/neuron-extensions/auth/accounts v0.0.0-20200828181848-ded3451554c6
	github.com/neuronlabs/neuron-extensions/auth/authenticator v0.0.0-20200828181848-ded3451554c6
	github.com/neuronlabs/neuron-extensions/auth/jwt-tokener v0.0.0-20200828181848-ded3451554c6
	github.com/neuronlabs/neuron-extensions/codec/json v0.0.3
	github.com/neuronlabs/neuron-extensions/repository/postgres v0.0.0-20200828181848-ded3451554c6
	github.com/neuronlabs/neuron-extensions/server/http v0.0.3
	github.com/neuronlabs/neuron-extensions/server/http/api/authentication v0.0.0-20200828181848-ded3451554c6
	github.com/neuronlabs/neuron-extensions/store/memory v0.0.0-20200828181848-ded3451554c6
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
)
