package radiobrowser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchPathEncoding(t *testing.T) {
	for _, tc := range []struct{ term, encoded string }{
		{"BBC Radio 1", "BBC%20Radio%201"},
		{"Радио", "%D0%A0%D0%B0%D0%B4%D0%B8%D0%BE"},
		{"100% hits", "100%25%20hits"},
		{"a/b?c=d#e", "a%2Fb%3Fc=d%23e"},
		{"a%2Fb", "a%252Fb"},
	} {
		t.Run(tc.term, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if want := "/json/stations/byname/" + tc.encoded; r.RequestURI != want {
					t.Errorf("RequestURI = %q, want %q", r.RequestURI, want)
				}
				if want := "/json/stations/byname/" + tc.term; r.URL.Path != want {
					t.Errorf("decoded Path = %q, want %q", r.URL.Path, want)
				}
				_, _ = w.Write([]byte(`[]`))
			}))
			defer server.Close()
			if _, err := StationsSearch(context.Background(), server.URL, StationByName, tc.term, StationsOptions{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStationURLQueryRoundTrip(t *testing.T) {
	streamURL := "https://example.com/Радио live?token=a+b&path=a%2Fb#player"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/stations/byurl" || r.URL.Query().Get("url") != streamURL {
			t.Errorf("stream URL lost in request: %s", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()
	if _, err := StationsByURL(context.Background(), server.URL, streamURL, StationsOptions{}); err != nil {
		t.Fatal(err)
	}
}
