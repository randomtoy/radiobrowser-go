package radiobrowser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestStationCheckResults_OK_NoStationUUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/checks" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"checkuuid":"c1",
				"stationuuid":"s1",
				"source":"de1.api.radio-browser.info",
				"codec":"MP3",
				"bitrate":128,
				"hls":0,
				"ok":1,
				"timestamp_iso8601":"2020-01-06T15:34:58Z",
				"timestamp":"2020-01-06 15:34:58",
				"urlcache":"http://stream.example.com/mp3_128",
				"metainfo_overrides_database":0,
				"public":1,
				"name":"Station Name",
				"description":"Station Description",
				"tags":"rock,jazz",
				"countrycode":"NL",
				"countrysubdivisioncode":"NL-NH",
				"homepage":"http://example.com",
				"favicon":"http://example.com/favicon.ico",
				"loadbalancer":"lb1",
				"server_software":"Icecast",
				"sampling":44100,
				"timing_ms":495,
				"languagecodes":"en,nl",
				"ssl_error":0,
				"geo_lat":52.37,
				"geo_long":4.90
			}
		]`))
	}))
	defer ts.Close()

	got, err := StationCheckResults(context.Background(), ts.URL, "", ChecksOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	sc := got[0]

	if sc.CheckUUID != "c1" {
		t.Fatalf("unexpected checkuuid: %q", sc.CheckUUID)
	}
	if sc.StationUUID != "s1" {
		t.Fatalf("unexpected stationuuid: %q", sc.StationUUID)
	}
	if sc.Source != "de1.api.radio-browser.info" {
		t.Fatalf("unexpected source: %q", sc.Source)
	}
	if sc.Codec != "MP3" {
		t.Fatalf("unexpected codec: %q", sc.Codec)
	}
	if sc.Bitrate != 128 {
		t.Fatalf("unexpected bitrate: %d", sc.Bitrate)
	}
	if sc.HLS != 0 {
		t.Fatalf("unexpected hls: %d", sc.HLS)
	}
	if sc.OK != 1 {
		t.Fatalf("unexpected ok: %d", sc.OK)
	}
	if sc.TimestampISO8601 != "2020-01-06T15:34:58Z" {
		t.Fatalf("unexpected timestamp_iso8601: %q", sc.TimestampISO8601)
	}
	if sc.Timestamp != "2020-01-06 15:34:58" {
		t.Fatalf("unexpected timestamp: %q", sc.Timestamp)
	}
	if sc.URLCache != "http://stream.example.com/mp3_128" {
		t.Fatalf("unexpected urlcache: %q", sc.URLCache)
	}
	if sc.MetaInfoOverridesDatabase != 0 {
		t.Fatalf("unexpected metainfo_overrides_database: %d", sc.MetaInfoOverridesDatabase)
	}
	if sc.Public == nil || *sc.Public != 1 {
		t.Fatalf("unexpected public: %v", sc.Public)
	}
	if sc.Name == nil || *sc.Name != "Station Name" {
		t.Fatalf("unexpected name: %v", sc.Name)
	}
	if sc.Description == nil || *sc.Description != "Station Description" {
		t.Fatalf("unexpected description: %v", sc.Description)
	}
	if sc.Tags == nil || *sc.Tags != "rock,jazz" {
		t.Fatalf("unexpected tags: %v", sc.Tags)
	}
	if sc.CountryCode == nil || *sc.CountryCode != "NL" {
		t.Fatalf("unexpected countrycode: %v", sc.CountryCode)
	}
	if sc.CountrySubdivisionCode == nil || *sc.CountrySubdivisionCode != "NL-NH" {
		t.Fatalf("unexpected countrysubdivisioncode: %v", sc.CountrySubdivisionCode)
	}
	if sc.Homepage == nil || *sc.Homepage != "http://example.com" {
		t.Fatalf("unexpected homepage: %v", sc.Homepage)
	}
	if sc.Favicon == nil || *sc.Favicon != "http://example.com/favicon.ico" {
		t.Fatalf("unexpected favicon: %v", sc.Favicon)
	}
	if sc.LoadBalancer == nil || *sc.LoadBalancer != "lb1" {
		t.Fatalf("unexpected loadbalancer: %v", sc.LoadBalancer)
	}
	if sc.ServerSoftware == nil || *sc.ServerSoftware != "Icecast" {
		t.Fatalf("unexpected server_software: %v", sc.ServerSoftware)
	}
	if sc.Sampling == nil || *sc.Sampling != 44100 {
		t.Fatalf("unexpected sampling: %v", sc.Sampling)
	}
	if sc.TimingMS != 495 {
		t.Fatalf("unexpected timing_ms: %d", sc.TimingMS)
	}
	if sc.LanguageCodes == nil || *sc.LanguageCodes != "en,nl" {
		t.Fatalf("unexpected languagecodes: %v", sc.LanguageCodes)
	}
	if sc.SSLError != 0 {
		t.Fatalf("unexpected ssl_error: %d", sc.SSLError)
	}
	if sc.GeoLat == nil || *sc.GeoLat != 52.37 {
		t.Fatalf("unexpected geo_lat: %v", sc.GeoLat)
	}
	if sc.GeoLong == nil || *sc.GeoLong != 4.90 {
		t.Fatalf("unexpected geo_long: %v", sc.GeoLong)
	}
}

func TestStationCheckResults_OK_WithStationUUID_PathEscaped(t *testing.T) {
	stationUUID := "a/b?c=d"
	escaped := url.PathEscape(stationUUID)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/checks/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationCheckResults(context.Background(), ts.URL, stationUUID, ChecksOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationCheckResults_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("lastcheckuuid") != "c10" {
			t.Fatalf("unexpected lastcheckuuid: %q", q.Get("lastcheckuuid"))
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

	_, err := StationCheckResults(context.Background(), ts.URL, "", ChecksOptions{
		LastCheckUUID: "  c10  ",
		Seconds:       60,
		Limit:         5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationCheckResults_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationCheckResults(context.Background(), ts.URL, "", ChecksOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationCheckResults_InvalidBaseURL(t *testing.T) {
	_, err := StationCheckResults(context.Background(), "://bad-url", "", ChecksOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationCheckResults_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationCheckResults(context.Background(), ts.URL, "", ChecksOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationCheckResults_StationUUIDTrimmed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/checks/s1" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationCheckResults(context.Background(), ts.URL, "  s1  ", ChecksOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationCheckSteps_OK_QueryUUIDs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/json/checksteps" {
			t.Fatalf("unexpected path: %s", r.URL.EscapedPath())
		}

		q := r.URL.Query()
		if q.Get("uuids") != "u1,u2" {
			t.Fatalf("unexpected uuids: %q", q.Get("uuids"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"stepuuid":"st1",
				"parent_stepuuid":null,
				"checkuuid":"c1",
				"stationuuid":"u1",
				"url":"http://stream.example.com:8000/streamname",
				"urltype":"STREAM",
				"error":null,
				"creation_iso8601":"2021-04-13T20:11:19Z"
			}
		]`))
	}))
	defer ts.Close()

	got, err := StationCheckSteps(context.Background(), ts.URL, []string{"u1", "u2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
	if got[0].StepUUID != "st1" {
		t.Fatalf("unexpected stepuuid: %q", got[0].StepUUID)
	}
	if got[0].StationUUID != "u1" {
		t.Fatalf("unexpected stationuuid: %q", got[0].StationUUID)
	}
	if got[0].URLType == nil || *got[0].URLType != "STREAM" {
		t.Fatalf("unexpected urltype: %v", got[0].URLType)
	}
}

func TestStationCheckSteps_TrimsAndSkipsEmpty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("uuids") != "u1,u2" {
			t.Fatalf("unexpected uuids: %q", q.Get("uuids"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StationCheckSteps(context.Background(), ts.URL, []string{"  u1  ", "", "u2", "   "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationCheckSteps_EmptyUUIDs(t *testing.T) {
	_, err := StationCheckSteps(context.Background(), "http://example.com", nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "uuids is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationCheckSteps_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := StationCheckSteps(context.Background(), ts.URL, []string{"u1"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStationCheckSteps_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := StationCheckSteps(context.Background(), ts.URL, []string{"u1"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStationCheckSteps_InvalidBaseURL(t *testing.T) {
	_, err := StationCheckSteps(context.Background(), "://bad-url", []string{"u1"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
