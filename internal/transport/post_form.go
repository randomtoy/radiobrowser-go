package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func PostFormJSON(ctx context.Context, baseURL, path string, form url.Values, out any) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}

	u := base.ResolveReference(&url.URL{Path: path})

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
