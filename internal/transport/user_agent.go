package transport

import (
	"context"
	"strings"
)

// DefaultUserAgent identifies the library. Update the version when releasing.
const DefaultUserAgent = "radiobrowser-go/0.1.0"

type userAgentKey struct{}

func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, userAgentKey{}, strings.TrimSpace(userAgent))
}

func userAgent(ctx context.Context) string {
	if value, ok := ctx.Value(userAgentKey{}).(string); ok && value != "" {
		return value
	}
	return DefaultUserAgent
}
