package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var defaultClient = &http.Client{}

func GetJSON(ctx context.Context, baseURL, path string, q url.Values, out any) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}

	u := base.ResolveReference(&url.URL{Path: path})
	if len(q) > 0 {
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent(ctx))

	resp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
