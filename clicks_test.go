package radiobrowser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestStationClicks_OK_NoStationUUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/clicks" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"stationuuid":"s1",
				"clickuuid":"c1",
				"clicktimestamp_iso8601":"2020-01-11T23:14:58Z",
				"clicktimestamp":"2020-01-11 23:14:58"
			}
		]`))
	}))
	defer ts.Close()

	got, err := StationClicks(context.Background(), ts.URL, "", ClicksOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	if got[0].StationUUID != "s1" {
		t.Fatalf("unexpected stationuuid: %q", got[0].StationUUID)
	}
	if got[0].ClickUUID != "c1" {
		t.Fatalf("unexpected clickuuid: %q", got[0].ClickUUID)
	}
	if got[0].ClickTimestampISO8601 != "2020-01-11T23:14:58Z" {
		t.Fatalf("unexpected clicktimestamp_iso8601: %q", got[0].ClickTimestampISO8601)
	}
	if got[0].ClickTimestamp != "2020-01-11 23:14:58" {
		t.Fatalf("unexpected clicktimestamp: %q", got[0].ClickTimestamp)
	}
}

func TestStationClicks_OK_WithStationUUID_PathEscaped(t *testing.T) {
	stationUUID := "a/b?c=d"
	escaped := url.PathEscape(stationUUID)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/clicks/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationClicks(context.Background(), ts.URL, stationUUID, ClicksOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationClicks_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("lastclickuuid") != "c10" {
			t.Fatalf("unexpected lastclickuuid: %q", q.Get("lastclickuuid"))
		}
		if q.Get("seconds") != "60" {
			t.Fatalf("unexpected seconds: %q", q.Get("seconds"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationClicks(context.Background(), ts.URL, "", ClicksOptions{
		LastClickUUID: "  c10  ",
		Seconds:       60,
		Limit:         5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationClicks_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationClicks(context.Background(), ts.URL, "", ClicksOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationClicks_InvalidBaseURL(t *testing.T) {
	_, err := StationClicks(context.Background(), "://bad-url", "", ClicksOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationClicks_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationClicks(context.Background(), ts.URL, "", ClicksOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationClicks_StationUUIDTrimmed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/clicks/s1" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationClicks(context.Background(), ts.URL, "  s1  ", ClicksOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationClickCounter_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/url/s1" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	resp, err := StationClickCounter(context.Background(), ts.URL, "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.OK != true {
		t.Fatalf("unexpected ok: %v", resp.OK)
	}
}

func TestStationClickCounter_PathEscaped(t *testing.T) {
	uuid := "a/b?c=d"
	escaped := url.PathEscape(uuid)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/url/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	_, err := StationClickCounter(context.Background(), ts.URL, uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationClickCounter_EmptyUUID(t *testing.T) {
	_, err := StationClickCounter(context.Background(), "http://example.com", "   ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationClickCounter_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationClickCounter(context.Background(), ts.URL, "s1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationClickCounter_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":`))
	}))
	defer ts.Close()

	_, err := StationClickCounter(context.Background(), ts.URL, "s1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationClickCounter_InvalidBaseURL(t *testing.T) {
	_, err := StationClickCounter(context.Background(), "://bad-url", "s1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
