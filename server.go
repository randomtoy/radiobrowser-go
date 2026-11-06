package radiobrowser

import (
	"context"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

// Stats returns statistics for the selected Radio Browser API server.
func Stats(ctx context.Context, baseURL string) (*ServerStats, error) {
	var out ServerStats
	if err := transport.GetJSON(ctx, baseURL, "/json/stats", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Servers returns API mirrors from /json/servers. Prefer client-side DNS lookup
// of all.api.radio-browser.info when available, as recommended by Radio Browser.
// Unlike StreamingServers, these hosts serve the Radio Browser API itself.
func Servers(ctx context.Context, baseURL string) ([]ServerMirror, error) {
	var out []ServerMirror
	if err := transport.GetJSON(ctx, baseURL, "/json/servers", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Config returns the active configuration of the selected API server.
func Config(ctx context.Context, baseURL string) (*ServerConfig, error) {
	var out ServerConfig
	if err := transport.GetJSON(ctx, baseURL, "/json/config", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
