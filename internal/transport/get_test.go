package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetJSON_OK_Decode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/json/countries" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("unexpected Accept header: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"Netherlands","stationcount":"123"}]`))
	}))
	defer ts.Close()

	type Country struct {
		Name         string `json:"name"`
		StationCount int    `json:"stationcount,string"`
	}

	var out []Country
	err := GetJSON(context.Background(), ts.URL, "/json/countries", nil, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out))
	}
	if out[0].Name != "Netherlands" {
		t.Fatalf("expected name=%q, got %q", "Netherlands", out[0].Name)
	}
	if out[0].StationCount != 123 {
		t.Fatalf("expected stationcount=%d, got %d", 123, out[0].StationCount)
	}
}

func TestGetJSON_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("order"); got != "stationcount" {
			t.Fatalf("expected order=stationcount, got %q", got)
		}
		if got := q.Get("reverse"); got != "true" {
			t.Fatalf("expected reverse=true, got %q", got)
		}
		if got := q.Get("offset"); got != "10" {
			t.Fatalf("expected offset=10, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	q := url.Values{}
	q.Set("order", "stationcount")
	q.Set("reverse", "true")
	q.Set("offset", "10")

	var out []any
	err := GetJSON(context.Background(), ts.URL, "/json/countries", q, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetJSON_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	var out any
	err := GetJSON(context.Background(), ts.URL, "/json/countries", nil, &out)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetJSON_InvalidBaseURL(t *testing.T) {
	var out any
	err := GetJSON(context.Background(), "://bad-url", "/json/countries", nil, &out)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid base url") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetJSON_NilContext(t *testing.T) {
	var out any
	err := GetJSON(nil, "http://example.com", "/json/countries", nil, &out) //nolint:SA1012 // intentionally passing nil to test nil context handling
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "context is nil" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetJSON_OutNil(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	err := GetJSON(context.Background(), ts.URL, "/json/countries", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetJSON_HTTPClientError(t *testing.T) {
	old := defaultClient
	t.Cleanup(func() { defaultClient = old })

	defaultClient = &http.Client{
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network down")
		}),
	}

	var out any
	err := GetJSON(context.Background(), "http://example.com", "/json/countries", nil, &out)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	// http.Client wraps the transport error with method+url context.
	if !strings.Contains(err.Error(), "network down") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetJSON_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"bad_json":`)
	}))
	defer ts.Close()

	var out map[string]any
	err := GetJSON(context.Background(), ts.URL, "/json/countries", nil, &out)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
