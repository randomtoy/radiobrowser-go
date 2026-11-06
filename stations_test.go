package radiobrowser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestStations_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("order") != "votes" {
			t.Fatalf("unexpected order: %q", q.Get("order"))
		}
		if q.Get("reverse") != "true" {
			t.Fatalf("unexpected reverse: %q", q.Get("reverse"))
		}
		if q.Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", q.Get("hidebroken"))
		}
		if q.Get("offset") != "10" {
			t.Fatalf("unexpected offset: %q", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"changeuuid":"c1",
				"stationuuid":"s1",
				"serveruuid":"srv1",
				"name":"Test Station",
				"url":"http://example.com/stream",
				"url_resolved":"http://example.com/stream",
				"homepage":"http://example.com",
				"favicon":"",
				"tags":"rock",
				"country":"Netherlands",
				"countrycode":"NL",
				"state":"",
				"language":"english",
				"languagecodes":"en",
				"votes":1,
				"lastchangetime":"2020-01-01 00:00:00",
				"lastchangetime_iso8601":"2020-01-01T00:00:00Z",
				"codec":"MP3",
				"bitrate":128,
				"hls":0,
				"lastcheckok":1,
				"lastchecktime":"2020-01-01 00:00:00",
				"lastchecktime_iso8601":"2020-01-01T00:00:00Z",
				"lastcheckoktime":"2020-01-01 00:00:00",
				"lastcheckoktime_iso8601":"2020-01-01T00:00:00Z",
				"lastlocalchecktime":"2020-01-01 00:00:00",
				"clicktimestamp":"2020-01-01 00:00:00",
				"clicktimestamp_iso8601":"2020-01-01T00:00:00Z",
				"clickcount":0,
				"clicktrend":0,
				"geo_lat":0.0,
				"geo_long":0.0,
				"has_extended_info":false
			}
		]`))
	}))
	defer ts.Close()

	got, err := Stations(context.Background(), ts.URL, StationsOptions{
		Order:      StationOrderVotes,
		Reverse:    true,
		HideBroken: true,
		Offset:     10,
		Limit:      5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
	if got[0].StationUUID != "s1" {
		t.Fatalf("unexpected stationuuid: %q", got[0].StationUUID)
	}
}

func TestStations_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := Stations(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStations_InvalidBaseURL(t *testing.T) {
	_, err := Stations(context.Background(), "://bad-url", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStations_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := Stations(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsSearch_OK_PathAndQuery(t *testing.T) {
	term := "new york"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/byname/new%20york" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("limit") != "3" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsSearch(context.Background(), ts.URL, StationByName, term, StationsOptions{Limit: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsSearch_EmptyKind(t *testing.T) {
	_, err := StationsSearch(context.Background(), "http://example.com", "", "abc", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "search kind is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsSearch_EmptyTerm(t *testing.T) {
	_, err := StationsSearch(context.Background(), "http://example.com", StationByName, "   ", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "search term is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsSearch_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsSearch(context.Background(), ts.URL, StationByName, "abc", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsSearch_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsSearch(context.Background(), ts.URL, StationByName, "abc", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByUUID_OK_QueryUUIDs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/byuuid" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		if r.URL.Query().Get("uuids") != "u1,u2" {
			t.Fatalf("unexpected uuids: %q", r.URL.Query().Get("uuids"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"changeuuid":"c1",
				"stationuuid":"u1",
				"serveruuid":"srv1",
				"name":"Test Station",
				"url":"http://example.com/stream",
				"url_resolved":"http://example.com/stream",
				"homepage":"http://example.com",
				"favicon":"",
				"tags":"rock",
				"country":"Netherlands",
				"countrycode":"NL",
				"state":"",
				"language":"english",
				"languagecodes":"en",
				"votes":1,
				"lastchangetime":"2020-01-01 00:00:00",
				"lastchangetime_iso8601":"2020-01-01T00:00:00Z",
				"codec":"MP3",
				"bitrate":128,
				"hls":0,
				"lastcheckok":1,
				"lastchecktime":"2020-01-01 00:00:00",
				"lastchecktime_iso8601":"2020-01-01T00:00:00Z",
				"lastcheckoktime":"2020-01-01 00:00:00",
				"lastcheckoktime_iso8601":"2020-01-01T00:00:00Z",
				"lastlocalchecktime":"2020-01-01 00:00:00",
				"clicktimestamp":"2020-01-01 00:00:00",
				"clicktimestamp_iso8601":"2020-01-01T00:00:00Z",
				"clickcount":0,
				"clicktrend":0,
				"geo_lat":0.0,
				"geo_long":0.0,
				"has_extended_info":false
			}
		]`))
	}))
	defer ts.Close()

	got, err := StationsByUUID(context.Background(), ts.URL, []string{"u1", "u2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
	if got[0].StationUUID != "u1" {
		t.Fatalf("unexpected stationuuid: %q", got[0].StationUUID)
	}
}

func TestStationsByUUID_TrimsAndSkipsEmpty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("uuids") != "u1,u2" {
			t.Fatalf("unexpected uuids: %q", r.URL.Query().Get("uuids"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsByUUID(context.Background(), ts.URL, []string{"  u1  ", "", "u2", "   "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByUUID_EmptyUUIDs(t *testing.T) {
	_, err := StationsByUUID(context.Background(), "http://example.com", nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "uuids is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByUUID_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsByUUID(context.Background(), ts.URL, []string{"u1"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByUUID_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsByUUID(context.Background(), ts.URL, []string{"u1"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByUUID_InvalidBaseURL(t *testing.T) {
	_, err := StationsByUUID(context.Background(), "://bad-url", []string{"u1"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsSearch_PathEscapes(t *testing.T) {
	term := "a/b?c=d"
	escaped := url.PathEscape(term)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/byname/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsSearch(context.Background(), ts.URL, StationByName, term, StationsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByURL_OK_URLQuery(t *testing.T) {
	rawURL := "http://example.com/a?b=c&d=e"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/byurl" || r.URL.Query().Get("url") != rawURL {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		if r.URL.Query().Get("limit") != "3" {
			t.Fatalf("unexpected limit: %q", r.URL.Query().Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsByURL(context.Background(), ts.URL, rawURL, StationsOptions{Limit: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByURL_EmptyURL(t *testing.T) {
	_, err := StationsByURL(context.Background(), "http://example.com", "   ", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "url is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByURL_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsByURL(context.Background(), ts.URL, "http://example.com/stream", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByURL_DecodeError(t *testing.T) {
	rawURL := "http://example.com/stream"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/byurl" || r.URL.Query().Get("url") != rawURL {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsByURL(context.Background(), ts.URL, rawURL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByURL_InvalidBaseURL(t *testing.T) {
	_, err := StationsByURL(context.Background(), "://bad-url", "http://example.com/stream", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByClicks_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/topclick" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("order") != "clickcount" {
			t.Fatalf("unexpected order: %q", q.Get("order"))
		}
		if q.Get("reverse") != "true" {
			t.Fatalf("unexpected reverse: %q", q.Get("reverse"))
		}
		if q.Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", q.Get("hidebroken"))
		}
		if q.Get("offset") != "10" {
			t.Fatalf("unexpected offset: %q", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsByClicks(context.Background(), ts.URL, StationsOptions{
		Order:      StationOrderClickCount,
		Reverse:    true,
		HideBroken: true,
		Offset:     10,
		Limit:      5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByClicks_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsByClicks(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByClicks_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsByClicks(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByClicks_InvalidBaseURL(t *testing.T) {
	_, err := StationsByClicks(context.Background(), "://bad-url", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByVotes_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/topvote" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("offset") != "10" {
			t.Fatalf("unexpected offset: %q", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}
		if q.Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", q.Get("hidebroken"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsByVotes(context.Background(), ts.URL, StationsOptions{
		Offset:     10,
		Limit:      5,
		HideBroken: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByVotes_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsByVotes(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByVotes_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsByVotes(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByVotes_InvalidBaseURL(t *testing.T) {
	_, err := StationsByVotes(context.Background(), "://bad-url", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByRecentClick_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/lastclick" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("offset") != "10" {
			t.Fatalf("unexpected offset: %q", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}
		if q.Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", q.Get("hidebroken"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsByRecentClick(context.Background(), ts.URL, StationsOptions{
		Offset:     10,
		Limit:      5,
		HideBroken: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByRecentClick_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsByRecentClick(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByRecentClick_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsByRecentClick(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByRecentClick_InvalidBaseURL(t *testing.T) {
	_, err := StationsByRecentClick(context.Background(), "://bad-url", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByRecentlyChanged_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/lastchange" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("offset") != "10" {
			t.Fatalf("unexpected offset: %q", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}
		if q.Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", q.Get("hidebroken"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsByRecentlyChanged(context.Background(), ts.URL, StationsOptions{
		Offset:     10,
		Limit:      5,
		HideBroken: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByRecentlyChanged_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsByRecentlyChanged(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsByRecentlyChanged_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsByRecentlyChanged(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsByRecentlyChanged_InvalidBaseURL(t *testing.T) {
	_, err := StationsByRecentlyChanged(context.Background(), "://bad-url", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsChanged_OK_NoStationUUID_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/changed" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("lastchangeuuid") != "lc1" {
			t.Fatalf("unexpected lastchangeuuid: %q", q.Get("lastchangeuuid"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsChanged(context.Background(), ts.URL, "", StationsChangedOptions{
		LastChangeUUID: "  lc1  ",
		Limit:          10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsChanged_OK_WithStationUUID_PathEscaped(t *testing.T) {
	uuid := "a/b?c=d"
	escaped := url.PathEscape(uuid)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/changed/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsChanged(context.Background(), ts.URL, uuid, StationsChangedOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsChanged_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsChanged(context.Background(), ts.URL, "", StationsChangedOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsChanged_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsChanged(context.Background(), ts.URL, "", StationsChangedOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsChanged_InvalidBaseURL(t *testing.T) {
	_, err := StationsChanged(context.Background(), "://bad-url", "", StationsChangedOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsBroken_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/broken" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("offset") != "10" {
			t.Fatalf("unexpected offset: %q", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}
		if q.Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", q.Get("hidebroken"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationsBroken(context.Background(), ts.URL, StationsOptions{
		Offset:     10,
		Limit:      5,
		HideBroken: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsBroken_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationsBroken(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationsBroken_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationsBroken(context.Background(), ts.URL, StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationsBroken_InvalidBaseURL(t *testing.T) {
	_, err := StationsBroken(context.Background(), "://bad-url", StationsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationSearch_OK_QueryParams(t *testing.T) {
	hasGeo := true
	isHTTPS := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/stations/search" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("name") != "Jazz FM" {
			t.Fatalf("unexpected name: %q", q.Get("name"))
		}
		if q.Get("nameExact") != "true" {
			t.Fatalf("unexpected nameExact: %q", q.Get("nameExact"))
		}
		if q.Get("country") != "Netherlands" {
			t.Fatalf("unexpected country: %q", q.Get("country"))
		}
		if q.Get("countryExact") != "true" {
			t.Fatalf("unexpected countryExact: %q", q.Get("countryExact"))
		}
		if q.Get("countrycode") != "NL" {
			t.Fatalf("unexpected countrycode: %q", q.Get("countrycode"))
		}
		if q.Get("state") != "Noord-Holland" {
			t.Fatalf("unexpected state: %q", q.Get("state"))
		}
		if q.Get("stateExact") != "true" {
			t.Fatalf("unexpected stateExact: %q", q.Get("stateExact"))
		}
		if q.Get("language") != "english" {
			t.Fatalf("unexpected language: %q", q.Get("language"))
		}
		if q.Get("languageExact") != "true" {
			t.Fatalf("unexpected languageExact: %q", q.Get("languageExact"))
		}
		if q.Get("tag") != "jazz" {
			t.Fatalf("unexpected tag: %q", q.Get("tag"))
		}
		if q.Get("tagExact") != "true" {
			t.Fatalf("unexpected tagExact: %q", q.Get("tagExact"))
		}
		if q.Get("tagList") != "jazz,blues" {
			t.Fatalf("unexpected tagList: %q", q.Get("tagList"))
		}
		if q.Get("codec") != "MP3" {
			t.Fatalf("unexpected codec: %q", q.Get("codec"))
		}
		if q.Get("bitrateMin") != "64" {
			t.Fatalf("unexpected bitrateMin: %q", q.Get("bitrateMin"))
		}
		if q.Get("bitrateMax") != "320" {
			t.Fatalf("unexpected bitrateMax: %q", q.Get("bitrateMax"))
		}
		if q.Get("has_geo_info") != "true" {
			t.Fatalf("unexpected has_geo_info: %q", q.Get("has_geo_info"))
		}
		if q.Get("is_https") != "false" {
			t.Fatalf("unexpected is_https: %q", q.Get("is_https"))
		}
		if q.Get("order") != "votes" {
			t.Fatalf("unexpected order: %q", q.Get("order"))
		}
		if q.Get("reverse") != "true" {
			t.Fatalf("unexpected reverse: %q", q.Get("reverse"))
		}
		if q.Get("offset") != "10" {
			t.Fatalf("unexpected offset: %q", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("unexpected limit: %q", q.Get("limit"))
		}
		if q.Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", q.Get("hidebroken"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationSearch(context.Background(), ts.URL, StationSearchOptions{
		Name:          "  Jazz FM  ",
		NameExact:     true,
		Country:       "Netherlands",
		CountryExact:  true,
		CountryCode:   "NL",
		State:         "Noord-Holland",
		StateExact:    true,
		Language:      "english",
		LanguageExact: true,
		Tag:           "jazz",
		TagExact:      true,
		TagList:       []string{" jazz ", "", "blues"},
		Codec:         "MP3",
		BitrateMin:    64,
		BitrateMax:    320,
		HasGeoInfo:    &hasGeo,
		IsHTTPS:       &isHTTPS,
		Order:         StationOrderVotes,
		Reverse:       true,
		Offset:        10,
		Limit:         5,
		HideBroken:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationSearch_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationSearch(context.Background(), ts.URL, StationSearchOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationSearch_InvalidBaseURL(t *testing.T) {
	_, err := StationSearch(context.Background(), "://bad-url", StationSearchOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationSearch_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationSearch(context.Background(), ts.URL, StationSearchOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
