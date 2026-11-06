package radiobrowser

import (
	"context"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

// DefaultUserAgent is sent when no application User-Agent is specified.
const DefaultUserAgent = transport.DefaultUserAgent

// WithUserAgent returns a context that sets the User-Agent for any API call made
// with it or a derived context. Use an application name and version, for example
// "MyRadio/1.2". Leading and trailing whitespace is removed; an empty value
// restores DefaultUserAgent. The parent context is not changed.
func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return transport.WithUserAgent(ctx, userAgent)
}
