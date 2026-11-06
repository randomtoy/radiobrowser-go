package radiobrowser

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

// Stations lists stations with optional sorting and pagination.
func Stations(ctx context.Context, baseURL string, opts StationsOptions) ([]Station, error) {
	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsSearch performs a simple substring or exact search selected by kind.
func StationsSearch(
	ctx context.Context,
	baseURL string,
	kind StationSearchKind,
	searchTerm string,
	opts StationsOptions,
) ([]Station, error) {
	searchTerm = strings.TrimSpace(searchTerm)
	if searchTerm == "" {
		return nil, fmt.Errorf("search term is empty")
	}
	if kind == "" {
		return nil, fmt.Errorf("search kind is empty")
	}

	path := "/json/stations/" + string(kind) + "/" + url.PathEscape(searchTerm)

	var out []Station
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsByUUID retrieves stations matching any of the supplied UUIDs.
func StationsByUUID(ctx context.Context, baseURL string, uuids []string) ([]Station, error) {
	items := make([]string, 0, len(uuids))
	for _, s := range uuids {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		items = append(items, s)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("uuids is empty")
	}

	q := url.Values{}
	q.Set("uuids", strings.Join(items, ","))

	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/byuuid", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsByURL finds stations whose stream URL exactly matches rawURL.
func StationsByURL(
	ctx context.Context,
	baseURL string,
	rawURL string,
	opts StationsOptions,
) ([]Station, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, fmt.Errorf("url is empty")
	}

	q := opts.toQuery()
	q.Set("url", rawURL)

	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/byurl", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsByClicks lists stations with the most clicks in the last 24 hours.
func StationsByClicks(ctx context.Context, baseURL string, opts StationsOptions) ([]Station, error) {
	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/topclick", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsByVotes lists stations with the most votes.
func StationsByVotes(ctx context.Context, baseURL string, opts StationsOptions) ([]Station, error) {
	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/topvote", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsByRecentClick lists stations in order of recent playback.
func StationsByRecentClick(ctx context.Context, baseURL string, opts StationsOptions) ([]Station, error) {
	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/lastclick", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsByRecentlyChanged lists recently added or changed stations.
func StationsByRecentlyChanged(ctx context.Context, baseURL string, opts StationsOptions) ([]Station, error) {
	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/lastchange", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsChanged retrieves old station versions; an empty UUID selects all stations.
func StationsChanged(
	ctx context.Context,
	baseURL string,
	stationUUID string,
	opts StationsChangedOptions,
) ([]Station, error) {
	path := "/json/stations/changed"

	stationUUID = strings.TrimSpace(stationUUID)
	if stationUUID != "" {
		path += "/" + url.PathEscape(stationUUID)
	}

	var out []Station
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationsBroken lists stations that failed their connection check.
func StationsBroken(ctx context.Context, baseURL string, opts StationsOptions) ([]Station, error) {
	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/broken", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StationSearch searches using combined filters, including optional geographic bounds.
func StationSearch(
	ctx context.Context,
	baseURL string,
	opts StationSearchOptions,
) ([]Station, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	var out []Station
	if err := transport.GetJSON(ctx, baseURL, "/json/stations/search", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}
