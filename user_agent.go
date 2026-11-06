package radiobrowser

import (
	"context"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

// DefaultUserAgent returns the library name and the module version recorded in
// the binary's build information. Unversioned builds and local module replacements
// use "radiobrowser-go/dev". The result is computed once at package initialization.
func DefaultUserAgent() string {
	return transport.DefaultUserAgent()
}

// WithUserAgent returns a context that sets the User-Agent for any API call made
// with it or a derived context. Use an application name and version, for example
// "MyRadio/1.2". Leading and trailing whitespace is removed; an empty value
// restores DefaultUserAgent(). The parent context is not changed.
func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return transport.WithUserAgent(ctx, userAgent)
}
