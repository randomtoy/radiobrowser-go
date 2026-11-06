package radiobrowser

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLimitOmissionAndExplicitZero(t *testing.T) {
	ctx := context.Background()
	calls := []struct {
		name string
		run  func(string, int, bool) error
	}{
		{"stations", func(base string, limit int, set bool) error {
			opts := StationsOptions{Limit: limit, LimitSet: set}
			_, err := Stations(ctx, base, opts)
			return err
		}},
		{"search", func(base string, limit int, set bool) error {
			opts := StationSearchOptions{Limit: limit, LimitSet: set}
			_, err := StationSearch(ctx, base, opts)
			return err
		}},
		{"changes", func(base string, limit int, set bool) error {
			opts := StationsChangedOptions{Limit: limit, LimitSet: set}
			_, err := StationsChanged(ctx, base, "", opts)
			return err
		}},
		{"checks", func(base string, limit int, set bool) error {
			opts := ChecksOptions{Limit: limit, LimitSet: set}
			_, err := StationCheckResults(ctx, base, "", opts)
			return err
		}},
		{"clicks", func(base string, limit int, set bool) error {
			opts := ClicksOptions{Limit: limit, LimitSet: set}
			_, err := StationClicks(ctx, base, "", opts)
			return err
		}},
		{"countries", func(base string, limit int, set bool) error {
			opts := CountriesOptions{Limit: limit, LimitSet: set}
			_, err := Countries(ctx, base, "", opts)
			return err
		}},
		{"codecs", func(base string, limit int, set bool) error {
			opts := CodecsOptions{Limit: limit, LimitSet: set}
			_, err := Codecs(ctx, base, "", opts)
			return err
		}},
		{"states", func(base string, limit int, set bool) error {
			opts := StatesOptions{Limit: limit, LimitSet: set}
			_, err := States(ctx, base, "", "", opts)
			return err
		}},
		{"languages", func(base string, limit int, set bool) error {
			opts := LanguagesOptions{Limit: limit, LimitSet: set}
			_, err := Languages(ctx, base, "", opts)
			return err
		}},
		{"tags", func(base string, limit int, set bool) error {
			opts := TagsOptions{Limit: limit, LimitSet: set}
			_, err := Tags(ctx, base, "", opts)
			return err
		}},
	}
	for _, call := range calls {
		t.Run(call.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				limit   int
				set     bool
				want    string
				present bool
			}{
				{"default", 0, false, "", false}, {"zero", 0, true, "0", true}, {"positive", 7, false, "7", true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						q := r.URL.Query()
						if q.Has("limit") != tc.present || q.Get("limit") != tc.want {
							t.Errorf("unexpected limit in %s", r.RequestURI)
						}
						_, _ = w.Write([]byte(`[]`))
					}))
					defer server.Close()
					if err := call.run(server.URL, tc.limit, tc.set); err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
}
func TestSearchGeographyAndTriStateOptions(t *testing.T) {
	zero := 0.0
	no := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for key, want := range map[string]string{"geo_lat": "0", "geo_long": "0", "geo_distance": "0", "has_geo_info": "false", "has_extended_info": "false", "is_https": "false"} {
			if !q.Has(key) || q.Get(key) != want {
				t.Errorf("%s = %q, want %q", key, q.Get(key), want)
			}
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()
	_, err := StationSearch(context.Background(), server.URL, StationSearchOptions{GeoLat: &zero, GeoLong: &zero, GeoDistance: &zero, HasGeoInfo: &no, HasExtendedInfo: &no, IsHTTPS: &no})
	if err != nil {
		t.Fatal(err)
	}
	q := (StationSearchOptions{}).toQuery()
	for _, key := range []string{"geo_lat", "geo_long", "geo_distance", "has_geo_info", "has_extended_info", "is_https"} {
		if q.Has(key) {
			t.Errorf("unset %s was sent", key)
		}
	}
}
func TestGeographicValidation(t *testing.T) {
	for _, tc := range []struct {
		field string
		value float64
	}{
		{"geo_lat", -90.1}, {"geo_lat", 90.1}, {"geo_lat", math.NaN()},
		{"geo_long", -180.1}, {"geo_long", 180.1}, {"geo_long", math.Inf(1)},
		{"geo_distance", -1}, {"geo_distance", math.Inf(1)}, {"geo_distance", math.NaN()},
	} {
		t.Run(tc.field, func(t *testing.T) {
			opts := StationSearchOptions{}
			switch tc.field {
			case "geo_lat":
				opts.GeoLat = &tc.value
			case "geo_long":
				opts.GeoLong = &tc.value
			case "geo_distance":
				opts.GeoDistance = &tc.value
			}
			_, err := StationSearch(context.Background(), "://invalid", opts)
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("expected geographic validation error, got %v", err)
			}
			if tc.field != "geo_distance" {
				_, err = AddStation(context.Background(), "://invalid", AddStationRequest{Name: "Test", URL: "https://example.com", GeoLat: opts.GeoLat, GeoLong: opts.GeoLong})
				if err == nil || !strings.Contains(err.Error(), tc.field) {
					t.Fatalf("expected add geographic validation error, got %v", err)
				}
			}
		})
	}
	for _, value := range []float64{-90, 0, 90} {
		if err := validateCoordinates(&value, nil); err != nil {
			t.Fatal(err)
		}
	}
	for _, value := range []float64{-180, 0, 180} {
		if err := validateCoordinates(nil, &value); err != nil {
			t.Fatal(err)
		}
	}
}
func TestAddStationModernFields(t *testing.T) {
	zero := 0.0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		for key, want := range map[string]string{"iso_3166_2": "NL-NH", "languagecodes": "dut,eng", "geo_lat": "0", "geo_long": "0"} {
			if r.Form.Get(key) != want {
				t.Errorf("%s = %q, want %q", key, r.Form.Get(key), want)
			}
		}
		if r.Form.Has("state") || r.Form.Has("language") {
			t.Error("legacy fields should be omitted when unset")
		}
		_, _ = w.Write([]byte(`{"ok":true,"uuid":"s1"}`))
	}))
	defer server.Close()
	_, err := AddStation(context.Background(), server.URL, AddStationRequest{Name: "Radio", URL: "https://example.com", ISO3166_2: " NL-NH ", LanguageCodes: " dut,eng ", GeoLat: &zero, GeoLong: &zero})
	if err != nil {
		t.Fatal(err)
	}
}
