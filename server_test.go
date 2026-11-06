package radiobrowser

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestServerJSONEndpoints(t *testing.T) {
	for _, tc := range []struct {
		name, path, payload string
		fetch               func(context.Context, string) (any, error)
		want                any
	}{
		{"stats", "/json/stats", `{"supported_version":1,"software_version":"0.7.45","status":"OK","stations":100,"stations_broken":2,"tags":30,"clicks_last_hour":4,"clicks_last_day":50,"languages":6,"countries":7}`,
			func(ctx context.Context, base string) (any, error) { return Stats(ctx, base) },
			&ServerStats{SupportedVersion: 1, SoftwareVersion: "0.7.45", Status: "OK", Stations: 100, StationsBroken: 2, Tags: 30, ClicksLastHour: 4, ClicksLastDay: 50, Languages: 6, Countries: 7}},
		{"mirrors", "/json/servers", `[{"ip":"192.0.2.1","name":"mirror.example"},{"ip":"2001:db8::1","name":"mirror.example"}]`,
			func(ctx context.Context, base string) (any, error) { return Servers(ctx, base) },
			[]ServerMirror{{IP: "192.0.2.1", Name: "mirror.example"}, {IP: "2001:db8::1", Name: "mirror.example"}}},
		{"config", "/json/config", `{"check_enabled":true,"prometheus_exporter_enabled":true,"pull_servers":["https://mirror.example"],"tcp_timeout_seconds":10,"broken_stations_never_working_timeout_seconds":20,"broken_stations_timeout_seconds":30,"checks_timeout_seconds":40,"click_valid_timeout_seconds":50,"clicks_timeout_seconds":60,"mirror_pull_interval_seconds":70,"update_caches_interval_seconds":80,"server_name":"test","server_location":"Earth","server_country_code":"NL","check_retries":2,"check_batchsize":100,"check_pause_seconds":3,"api_threads":4,"cache_type":"redis","cache_ttl":120,"language_replace_filepath":"https://example.com/replace.csv","language_to_code_filepath":"https://example.com/codes.csv"}`,
			func(ctx context.Context, base string) (any, error) { return Config(ctx, base) },
			&ServerConfig{CheckEnabled: true, PrometheusExporterEnabled: true, PullServers: []string{"https://mirror.example"}, TCPTimeoutSeconds: 10, BrokenStationsNeverWorkingTimeoutSeconds: 20, BrokenStationsTimeoutSeconds: 30, ChecksTimeoutSeconds: 40, ClickValidTimeoutSeconds: 50, ClicksTimeoutSeconds: 60, MirrorPullIntervalSeconds: 70, UpdateCachesIntervalSeconds: 80, ServerName: "test", ServerLocation: "Earth", ServerCountryCode: "NL", CheckRetries: 2, CheckBatchSize: 100, CheckPauseSeconds: 3, APIThreads: 4, CacheType: "redis", CacheTTL: 120, LanguageReplaceFilepath: "https://example.com/replace.csv", LanguageToCodeFilepath: "https://example.com/codes.csv"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodGet || r.RequestURI != tc.path || r.Header.Get("User-Agent") != "TestApp/1" {
						t.Errorf("unexpected request: %s %s, UA %q", r.Method, r.RequestURI, r.UserAgent())
					}
					_, _ = w.Write([]byte(tc.payload))
				}))
				defer server.Close()
				got, err := tc.fetch(WithUserAgent(context.Background(), "TestApp/1"), server.URL)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("got %+v, want %+v", got, tc.want)
				}
			})
			for _, status := range []int{http.StatusOK, http.StatusServiceUnavailable} {
				t.Run(http.StatusText(status), func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(status)
						_, _ = w.Write([]byte(`{"invalid":`))
					}))
					defer server.Close()
					if _, err := tc.fetch(context.Background(), server.URL); err == nil {
						t.Fatal("expected HTTP or decoding error")
					}
				})
			}
			if _, err := tc.fetch(nil, "https://example.com"); err == nil || !strings.Contains(err.Error(), "context is nil") {
				t.Fatalf("nil context error: %v", err)
			}
			if _, err := tc.fetch(context.Background(), "://invalid"); err == nil {
				t.Fatal("expected invalid URL error")
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			if _, err := tc.fetch(ctx, "https://example.com"); !errors.Is(err, context.Canceled) {
				t.Fatalf("cancellation error: %v", err)
			}
		})
	}
}
