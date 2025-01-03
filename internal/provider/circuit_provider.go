package provider

import (
	"context"
	"github.com/sony/gobreaker/v2"
	"open_url_service/pkg/config"
)

type circuitBreaker struct {
	cb *gobreaker.CircuitBreaker[[]byte]
}

func (c circuitBreaker) Execute(ctx context.Context, fn ExecuteServiceFunc, options httpclient.RequestOptions) (any, error) {
	return c.cb.Execute(func() (interface{}, error) {
		resp, err := fn(options)
		return resp, err
	})
}

func (c circuitBreaker) GetState(ctx context.Context) gobreaker.State {
	return c.cb.State()
}

func NewCircuitBreaker(cfg *config.Config, settings gobreaker.Settings) CircuitBreaker {
	cb := gobreaker.NewCircuitBreaker(settings)
	return &circuitBreaker{cb: cb}
}
