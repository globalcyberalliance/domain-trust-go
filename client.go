package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fxamacker/cbor/v2"
)

const (
	DefaultTimeout = 30 * time.Second
	DocsURL        = "https://domain-trust.docs.globalcyberalliance.org"
	EndpointURL    = "https://domain-trust.globalcyberalliance.org/api"
	Version        = "2.0.0"
)

// Client represents the Domain Trust API client.
type Client struct {
	apiKey       string
	client       *http.Client
	contentType  string
	debug        bool
	encodingType string
}

// New initializes a new Domain Trust API client using the provided API key and options.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:       apiKey,
		client:       &http.Client{Timeout: DefaultTimeout},
		contentType:  ContentTypeCBOR,
		debug:        false,
		encodingType: EncodingTypeZSTD,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetAPIKey replaces the existing API key in use.
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// SetTimeout updates the client timeout.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// Option is a function that applies a configuration option to a Client.
type Option func(*Client)

// WithClient allows providing a custom *http.Client.
func WithClient(client *http.Client) Option {
	return func(c *Client) {
		if client == nil {
			client = &http.Client{Timeout: DefaultTimeout}
		}

		c.client = client
	}
}

// WithContentType overrides the default content type from CBOR to a user-specified value.
func WithContentType(contentType string) Option {
	return func(c *Client) {
		if contentType == "" {
			contentType = ContentTypeCBOR
		}

		c.contentType = contentType
	}
}

// WithDebug enables or disables debug mode.
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithEncodingType overrides the default encoding type from ZSTD to a user-specified value.
func WithEncodingType(encodingType string) Option {
	return func(c *Client) {
		if encodingType == "" {
			encodingType = EncodingTypeZSTD
		}

		c.encodingType = encodingType
	}
}

// WithTimeout sets a custom timeout on the underlying http.Client.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.client == nil {
			c.client = &http.Client{}
		}

		c.client.Timeout = timeout
	}
}

func (c *Client) marshal(v any) ([]byte, error) {

	switch c.contentType {
	case ContentTypeCBOR:
		return cbor.Marshal(v)
	case ContentTypeJSON:
		return json.Marshal(v)
	}

	return nil, fmt.Errorf("unsupported content type: %s", c.contentType)
}
