package greeterendpoint

import (
	"context"

	"github.com/go-kit/kit/log"

	greeterservice "github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterservice"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a greeter service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	GetHealthEndpoint   endpoint.Endpoint // used by Consul for the healthcheck
	GetGreetingEndpoint endpoint.Endpoint
}

// MakeEndpoints returns service Endoints, and wires in all the provided
// middlewares.
func MakeEndpoints(s greeterservice.Service, logger log.Logger) Endpoints {
	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = MakeGetHealthEndpoint(s)
		healthEndpoint = LoggingMiddleware(log.With(logger, "method", "GetHealth"))(healthEndpoint)
	}

	var greeterEndpoint endpoint.Endpoint
	{
		greeterEndpoint = MakeGetGreetingEndpoint(s)
		greeterEndpoint = LoggingMiddleware(log.With(logger, "method", "GetGreeting"))(greeterEndpoint)
	}

	return Endpoints{
		GetHealthEndpoint:   healthEndpoint,
		GetGreetingEndpoint: greeterEndpoint,
	}
}

// MakeGetHealthEndpoint constructs a GetHealth endpoint wrapping the service.
func MakeGetHealthEndpoint(s greeterservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health, err := s.GetHealth()
		return getHealthResponse{Health: health, Err: err}, nil
	}
}

// MakeGetGreetingEndpoint constructs a GetGreeter endpoint wrapping the service.
func MakeGetGreetingEndpoint(s greeterservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getGreetingRequest)
		greeting, err := s.GetGreeting(ctx, req.Name)
		return getGreetingResponse{Greeting: greeting, Err: err}, nil
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

type getHealthRequest struct{}

type getHealthResponse struct {
	Health bool  `json:"health,omitempty"`
	Err    error `json:"err,omitempty"`
}

func (r getHealthResponse) Failed() error { return r.Err }

type getGreetingRequest struct {
	Name string
}

type getGreetingResponse struct {
	Greeting string `json:"greeting,omitempty"`
	Err      error  `json:"err,omitempty"`
}

func (r getGreetingResponse) Failed() error { return r.Err }
