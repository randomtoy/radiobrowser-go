package radiobrowser

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

type ChecksOptions struct {
	LastCheckUUID string
	Seconds       int
	Limit         int
	// LimitSet sends Limit even when it is zero.
	LimitSet bool
}

func (o ChecksOptions) toQuery() url.Values {
	q := url.Values{}

	if s := strings.TrimSpace(o.LastCheckUUID); s != "" {
		q.Set("lastcheckuuid", s)
	}
	if o.Seconds > 0 {
		q.Set("seconds", strconv.Itoa(o.Seconds))
	}
	if o.Limit > 0 || o.LimitSet {
		q.Set("limit", strconv.Itoa(o.Limit))
	}

	return q
}

func StationCheckResults(
	ctx context.Context,
	baseURL string,
	stationUUID string,
	opts ChecksOptions,
) ([]StationCheck, error) {
	path := "/json/checks"

	stationUUID = strings.TrimSpace(stationUUID)
	if stationUUID != "" {
		path += "/" + url.PathEscape(stationUUID)
	}

	var out []StationCheck
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func StationCheckSteps(ctx context.Context, baseURL string, stationUUIDs []string) ([]StationCheckStep, error) {
	uuids := make([]string, 0, len(stationUUIDs))
	for _, s := range stationUUIDs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		uuids = append(uuids, s)
	}

	if len(uuids) == 0 {
		return nil, fmt.Errorf("uuids is empty")
	}

	q := url.Values{}
	q.Set("uuids", strings.Join(uuids, ","))

	var out []StationCheckStep
	if err := transport.GetJSON(ctx, baseURL, "/json/checksteps", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
