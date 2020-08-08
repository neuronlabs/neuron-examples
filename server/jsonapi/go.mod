module github.com/neuronlabs/neuron-examples/server/jsonapi

replace (
	github.com/neuronlabs/neuron => ./../../../neuron
	github.com/neuronlabs/neuron-extensions/codec/jsonapi => ./../../../neuron-extensions/codec/jsonapi
	github.com/neuronlabs/neuron-extensions/repository/postgres => ./../../../neuron-extensions/repository/postgres
	github.com/neuronlabs/neuron-extensions/server/http => ./../../../neuron-extensions/server/http
	github.com/neuronlabs/neuron-extensions/server/http/api/jsonapi => ./../../../neuron-extensions/server/http/api/jsonapi
)

go 1.14

require (
	github.com/neuronlabs/neuron v0.15.0
	github.com/neuronlabs/neuron-extensions/repository/postgres v0.0.0-00010101000000-000000000000
	github.com/neuronlabs/neuron-extensions/server/http v0.0.0
	github.com/neuronlabs/neuron-extensions/server/http/api/jsonapi v0.0.0-00010101000000-000000000000
)
