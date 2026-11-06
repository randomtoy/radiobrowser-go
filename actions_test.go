package radiobrowser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAddStation_OK_FormParams(t *testing.T) {
	lat := 52.37
	lon := 4.90

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.EscapedPath() != "/json/add" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			t.Fatalf("unexpected content-type: %q", ct)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}

		assertForm := func(k, want string) {
			if got := r.Form.Get(k); got != want {
				t.Fatalf("unexpected %s: %q", k, got)
			}
		}

		assertForm("name", "My Station")
		assertForm("url", "http://example.com/stream.mp3")
		assertForm("homepage", "http://example.com")
		assertForm("favicon", "http://example.com/favicon.png")
		assertForm("countrycode", "NL")
		assertForm("state", "Noord-Holland")
		assertForm("language", "English")
		assertForm("tags", "jazz,rock")
		assertForm("geo_lat", "52.37")
		assertForm("geo_long", "4.9")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"message":"station was added","uuid":"u1"}`))
	}))
	defer ts.Close()

	resp, err := AddStation(context.Background(), ts.URL, AddStationRequest{
		Name:        "  My Station  ",
		URL:         "  http://example.com/stream.mp3 ",
		Homepage:    "http://example.com",
		Favicon:     "http://example.com/favicon.png",
		CountryCode: "NL",
		State:       "Noord-Holland",
		Language:    "English",
		Tags:        "jazz,rock",
		GeoLat:      &lat,
		GeoLong:     &lon,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.OK {
		t.Fatalf("unexpected ok: %v", resp.OK)
	}
	if resp.UUID != "u1" {
		t.Fatalf("unexpected uuid: %q", resp.UUID)
	}
	if resp.Message == "" {
		t.Fatalf("expected non-empty message")
	}
}

func TestAddStation_Validation_NameEmpty(t *testing.T) {
	_, err := AddStation(context.Background(), "http://example.com", AddStationRequest{
		Name: "   ",
		URL:  "http://example.com/stream",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "name is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddStation_Validation_URLEmpty(t *testing.T) {
	_, err := AddStation(context.Background(), "http://example.com", AddStationRequest{
		Name: "x",
		URL:  "   ",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "url is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddStation_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := AddStation(context.Background(), ts.URL, AddStationRequest{
		Name: "x",
		URL:  "http://example.com/stream",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddStation_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":`))
	}))
	defer ts.Close()

	_, err := AddStation(context.Background(), ts.URL, AddStationRequest{
		Name: "x",
		URL:  "http://example.com/stream",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestAddStation_InvalidBaseURL(t *testing.T) {
	_, err := AddStation(context.Background(), "://bad-url", AddStationRequest{
		Name: "x",
		URL:  "http://example.com/stream",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestAddStation_UsesPOST(t *testing.T) {
	var gotMethod string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"message":"station was added","uuid":"u1"}`))
	}))
	defer ts.Close()

	_, err := AddStation(context.Background(), ts.URL, AddStationRequest{
		Name: "x",
		URL:  "http://example.com/stream",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("unexpected method: %s", gotMethod)
	}
}

// Compile-time guard for unused import in this file when you paste into your repo.
var _ = url.Values{}

func TestVoteForStation_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/vote/s1" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"message":"voted for station successfully"}`))
	}))
	defer ts.Close()

	resp, err := VoteForStation(context.Background(), ts.URL, "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OK != true {
		t.Fatalf("unexpected ok: %v", resp.OK)
	}
	if resp.Message == "" {
		t.Fatalf("expected non-empty message")
	}
}

func TestVoteForStation_PathEscaped(t *testing.T) {
	uuid := "a/b?c=d"
	escaped := url.PathEscape(uuid)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/vote/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"message":"ok"}`))
	}))
	defer ts.Close()

	_, err := VoteForStation(context.Background(), ts.URL, uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVoteForStation_EmptyUUID(t *testing.T) {
	_, err := VoteForStation(context.Background(), "http://example.com", "   ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "stationuuid is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVoteForStation_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := VoteForStation(context.Background(), ts.URL, "s1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVoteForStation_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":`))
	}))
	defer ts.Close()

	_, err := VoteForStation(context.Background(), ts.URL, "s1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestVoteForStation_InvalidBaseURL(t *testing.T) {
	_, err := VoteForStation(context.Background(), "://bad-url", "s1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
