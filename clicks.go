package radiobrowser

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

type ClicksOptions struct {
	LastClickUUID string
	Seconds       int
	Limit         int
}

func (o ClicksOptions) toQuery() url.Values {
	q := url.Values{}

	if s := strings.TrimSpace(o.LastClickUUID); s != "" {
		q.Set("lastclickuuid", s)
	}
	if o.Seconds > 0 {
		q.Set("seconds", strconv.Itoa(o.Seconds))
	}
	if o.Limit > 0 {
		q.Set("limit", strconv.Itoa(o.Limit))
	}

	return q
}

func StationClicks(
	ctx context.Context,
	baseURL string,
	stationUUID string,
	opts ClicksOptions,
) ([]StationClick, error) {
	path := "/json/clicks"

	stationUUID = strings.TrimSpace(stationUUID)
	if stationUUID != "" {
		path += "/" + url.PathEscape(stationUUID)
	}

	var out []StationClick
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func StationClickCounter(
	ctx context.Context,
	baseURL string,
	stationUUID string,
) (*StationClickCounterResponse, error) {
	stationUUID = strings.TrimSpace(stationUUID)
	if stationUUID == "" {
		return nil, fmt.Errorf("stationuuid is empty")
	}

	path := "/json/url/" + url.PathEscape(stationUUID)

	var out StationClickCounterResponse
	if err := transport.GetJSON(ctx, baseURL, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
