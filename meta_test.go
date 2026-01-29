package radiobrowser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCountries_OK_NoFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/countries" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"Albania","iso_3166_1":"AL","stationcount":"1"},
			{"name":"Netherlands","iso_3166_1":"NL","stationcount":"123"}
		]`))
	}))
	defer ts.Close()

	got, err := Countries(context.Background(), ts.URL, "", CountriesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}

	if got[0].Name != "Albania" {
		t.Fatalf("unexpected name: %q", got[0].Name)
	}
	if got[0].ISO3166_1 != "AL" {
		t.Fatalf("unexpected iso code: %q", got[0].ISO3166_1)
	}
	if got[0].StationCount != 1 {
		t.Fatalf("unexpected stationcount: %d", got[0].StationCount)
	}
}

func TestCountries_OK_WithFilter_PathEscaped(t *testing.T) {
	filter := "new york"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/countries/new%20york" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Countries(context.Background(), ts.URL, filter, CountriesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCountries_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("order") != "stationcount" {
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Countries(context.Background(), ts.URL, "", CountriesOptions{
		Order:      "stationcount",
		Reverse:    true,
		HideBroken: true,
		Offset:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCountries_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := Countries(context.Background(), ts.URL, "", CountriesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCountries_InvalidBaseURL(t *testing.T) {
	_, err := Countries(context.Background(), "://bad-url", "", CountriesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCountries_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := Countries(context.Background(), ts.URL, "", CountriesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCountries_FilterTrimmed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/countries/abc" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Countries(context.Background(), ts.URL, "  abc  ", CountriesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCountries_FilterSpecialCharsEscaped(t *testing.T) {
	filter := "a/b?c=d"
	escaped := url.PathEscape(filter)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/countries/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Countries(context.Background(), ts.URL, filter, CountriesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCodecs_OK_NoFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/codecs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"MP3","stationcount":"10"},
			{"name":"AAC","stationcount":"2"}
		]`))
	}))
	defer ts.Close()

	got, err := Codecs(context.Background(), ts.URL, "", CodecsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}

	if got[0].Name != "MP3" {
		t.Fatalf("unexpected name: %q", got[0].Name)
	}
	if got[0].StationCount != 10 {
		t.Fatalf("unexpected stationcount: %d", got[0].StationCount)
	}
}

func TestCodecs_OK_WithFilter_PathEscaped(t *testing.T) {
	filter := "a a"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/codecs/a%20a" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Codecs(context.Background(), ts.URL, filter, CodecsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCodecs_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("order") != "stationcount" {
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Codecs(context.Background(), ts.URL, "", CodecsOptions{
		Order:      "stationcount",
		Reverse:    true,
		HideBroken: true,
		Offset:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCodecs_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := Codecs(context.Background(), ts.URL, "", CodecsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCodecs_InvalidBaseURL(t *testing.T) {
	_, err := Codecs(context.Background(), "://bad-url", "", CodecsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCodecs_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := Codecs(context.Background(), ts.URL, "", CodecsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCodecs_FilterTrimmed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/codecs/abc" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Codecs(context.Background(), ts.URL, "  abc  ", CodecsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCodecs_FilterSpecialCharsEscaped(t *testing.T) {
	filter := "a/b?c=d"
	escaped := url.PathEscape(filter)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/codecs/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Codecs(context.Background(), ts.URL, filter, CodecsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStates_OK_NoCountry_NoFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/states" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"California","country":"United States","stationcount":"10"},
			{"name":"Bavaria","country":"Germany","stationcount":"2"}
		]`))
	}))
	defer ts.Close()

	got, err := States(context.Background(), ts.URL, "", "", StatesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].Name != "California" {
		t.Fatalf("unexpected name: %q", got[0].Name)
	}
	if got[0].Country != "United States" {
		t.Fatalf("unexpected country: %q", got[0].Country)
	}
	if got[0].StationCount != 10 {
		t.Fatalf("unexpected stationcount: %d", got[0].StationCount)
	}
}

func TestStates_OK_FilterOnly_PathEscaped(t *testing.T) {
	filter := "new york"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/states/new%20york" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, "", filter, StatesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStates_OK_CountryOnly_PathEscaped(t *testing.T) {
	country := "United States"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/states/United%20States" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, country, "", StatesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStates_OK_CountryAndFilter_PathEscaped(t *testing.T) {
	country := "United States"
	filter := "new york"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/states/United%20States/new%20york" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, country, filter, StatesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStates_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("order") != "stationcount" {
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, "", "", StatesOptions{
		Order:      "stationcount",
		Reverse:    true,
		HideBroken: true,
		Offset:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStates_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, "", "", StatesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStates_InvalidBaseURL(t *testing.T) {
	_, err := States(context.Background(), "://bad-url", "", "", StatesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStates_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, "", "", StatesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStates_TrimmedInputs(t *testing.T) {
	country := "  United States  "
	filter := "  new york  "

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/states/United%20States/new%20york" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, country, filter, StatesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStates_EscapesSpecialChars(t *testing.T) {
	country := "a/b?c=d"
	filter := "x y"
	wantCountry := url.PathEscape(country)
	wantFilter := url.PathEscape(filter)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/states/"+wantCountry+"/"+wantFilter {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := States(context.Background(), ts.URL, country, filter, StatesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLanguages_OK_NoFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/languages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"english","iso_639":"en","stationcount":"10"},
			{"name":"dutch","iso_639":"nl","stationcount":"2"}
		]`))
	}))
	defer ts.Close()

	got, err := Languages(context.Background(), ts.URL, "", LanguagesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}

	if got[0].Name != "english" {
		t.Fatalf("unexpected name: %q", got[0].Name)
	}
	if got[0].ISO639 != "en" {
		t.Fatalf("unexpected iso_639: %q", got[0].ISO639)
	}
	if got[0].StationCount != 10 {
		t.Fatalf("unexpected stationcount: %d", got[0].StationCount)
	}
}

func TestLanguages_OK_WithFilter_PathEscaped(t *testing.T) {
	filter := "new york"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/languages/new%20york" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Languages(context.Background(), ts.URL, filter, LanguagesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLanguages_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("order") != "stationcount" {
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Languages(context.Background(), ts.URL, "", LanguagesOptions{
		Order:      "stationcount",
		Reverse:    true,
		HideBroken: true,
		Offset:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLanguages_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := Languages(context.Background(), ts.URL, "", LanguagesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLanguages_InvalidBaseURL(t *testing.T) {
	_, err := Languages(context.Background(), "://bad-url", "", LanguagesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLanguages_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := Languages(context.Background(), ts.URL, "", LanguagesOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLanguages_FilterTrimmed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/languages/abc" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Languages(context.Background(), ts.URL, "  abc  ", LanguagesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLanguages_FilterSpecialCharsEscaped(t *testing.T) {
	filter := "a/b?c=d"
	escaped := url.PathEscape(filter)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/languages/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Languages(context.Background(), ts.URL, filter, LanguagesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLanguages_OK_MissingISO639(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"unknown","stationcount":"1"}
		]`))
	}))
	defer ts.Close()

	got, err := Languages(context.Background(), ts.URL, "", LanguagesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
	if got[0].ISO639 != "" {
		t.Fatalf("unexpected iso_639: %q", got[0].ISO639)
	}
}

func TestTags_OK_NoFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/tags" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"rock","stationcount":"10"},
			{"name":"jazz","stationcount":"2"}
		]`))
	}))
	defer ts.Close()

	got, err := Tags(context.Background(), ts.URL, "", TagsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}

	if got[0].Name != "rock" {
		t.Fatalf("unexpected name: %q", got[0].Name)
	}
	if got[0].StationCount != 10 {
		t.Fatalf("unexpected stationcount: %d", got[0].StationCount)
	}
}

func TestTags_OK_WithFilter_PathEscaped(t *testing.T) {
	filter := "new york"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/tags/new%20york" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Tags(context.Background(), ts.URL, filter, TagsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTags_OK_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("order") != "stationcount" {
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Tags(context.Background(), ts.URL, "", TagsOptions{
		Order:      "stationcount",
		Reverse:    true,
		HideBroken: true,
		Offset:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTags_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusTeapot)
	}))
	defer ts.Close()

	_, err := Tags(context.Background(), ts.URL, "", TagsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 418") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTags_InvalidBaseURL(t *testing.T) {
	_, err := Tags(context.Background(), "://bad-url", "", TagsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestTags_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad_json":`))
	}))
	defer ts.Close()

	_, err := Tags(context.Background(), ts.URL, "", TagsOptions{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestTags_FilterTrimmed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/tags/abc" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Tags(context.Background(), ts.URL, "  abc  ", TagsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTags_FilterSpecialCharsEscaped(t *testing.T) {
	filter := "a/b?c=d"
	escaped := url.PathEscape(filter)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/tags/"+escaped {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := Tags(context.Background(), ts.URL, filter, TagsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStreamingServers_OK_NoQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/streamingservers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"uuid":"u1",
				"url":"https://my.url/",
				"statusurl":"https://my.url/status.xsl",
				"admin":"tester@example.com",
				"location":"Mars, high street 3",
				"software":"Icecast 2.4.4"
			},
			{
				"uuid":"u2",
				"url":"http://my.url/",
				"error":"ResultDecodeError",
				"admin":"a",
				"location":"l",
				"software":"s"
			}
		]`))
	}))
	defer ts.Close()

	got, err := StreamingServers(context.Background(), ts.URL, StreamingServersOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}

	if got[0].UUID != "u1" {
		t.Fatalf("unexpected uuid: %q", got[0].UUID)
	}
	if got[0].StatusURL == nil || *got[0].StatusURL != "https://my.url/status.xsl" {
		t.Fatalf("unexpected statusurl: %v", got[0].StatusURL)
	}

	if got[1].UUID != "u2" {
		t.Fatalf("unexpected uuid: %q", got[1].UUID)
	}
	if got[1].Error == nil || *got[1].Error != "ResultDecodeError" {
		t.Fatalf("unexpected error: %v", got[1].Error)
	}
	if got[1].StatusURL != nil {
		t.Fatalf("unexpected statusurl: %v", got[1].StatusURL)
	}
}

func TestStreamingServers_OK_HideBroken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("hidebroken") != "true" {
			t.Fatalf("unexpected hidebroken: %q", r.URL.Query().Get("hidebroken"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := StreamingServers(context.Background(), ts.URL, StreamingServersOptions{HideBroken: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
