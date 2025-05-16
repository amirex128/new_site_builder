package v8

import (
	"net/http"
)

// ESConfig is a public config type for Elasticsearch connections, decoupled from the elasticsearch package.
type ESConfig struct {
	Addresses []string
	Username  string
	Password  string
	Transport http.RoundTripper
	Logger    Logger
}

// ESClient is a minimal interface for Elasticsearch operations you want to expose.
type ESClient interface {
	Info() (ESResponse, error)
	Search(index string, body string) (ESResponse, error)
	Index(index string, id string, body string) (ESResponse, error)
	Delete(index string, id string) (ESResponse, error)
	Update(index string, id string, body string) (ESResponse, error)
	Get(index string, id string) (ESResponse, error)
	Close() error
}

// ESResponse is a minimal response wrapper for Info and other methods.
type ESResponse interface {
	String() string
	Status() int
	Body() []byte
	IsError() bool
}
