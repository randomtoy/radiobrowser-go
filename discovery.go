package radiobrowser

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

// DiscoverServers returns unique HTTPS base URLs for Radio Browser API servers
// in randomized order. It queries _api._tcp.radio-browser.info through
// net.DefaultResolver, without requiring a known API server or making HTTP
// requests. Set a context deadline to bound the DNS lookup.
//
// Discovery does not check server health. Callers can use the first URL with any
// API function and try the remaining URLs if a request fails. An empty or unusable
// DNS result is returned as an error.
func DiscoverServers(ctx context.Context) ([]string, error) {
	return discoverServers(ctx, net.DefaultResolver.LookupSRV)
}

func discoverServers(ctx context.Context, lookup func(context.Context, string, string, string) (string, []*net.SRV, error)) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	_, records, err := lookup(ctx, "api", "tcp", "radio-browser.info")
	if err != nil {
		return nil, fmt.Errorf("discover Radio Browser servers: %w", err)
	}

	var servers []string
	seen := make(map[string]bool, len(records))
	for _, record := range records {
		if record == nil || record.Port == 0 {
			continue
		}
		host := strings.ToLower(strings.TrimSuffix(record.Target, "."))
		if host == "" {
			continue
		}
		if record.Port != 443 {
			host = net.JoinHostPort(host, strconv.Itoa(int(record.Port)))
		}
		baseURL := "https://" + host
		if !seen[baseURL] {
			seen[baseURL] = true
			servers = append(servers, baseURL)
		}
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no Radio Browser servers found in DNS")
	}

	rand.Shuffle(len(servers), func(i, j int) {
		servers[i], servers[j] = servers[j], servers[i]
	})
	return servers, nil
}
