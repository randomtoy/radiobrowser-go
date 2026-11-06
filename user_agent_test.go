package radiobrowser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	rb "github.com/randomtoy/radiobrowser-go"
)

func TestUserAgent(t *testing.T) {
	parent := context.Background()
	app := rb.WithUserAgent(parent, "MyRadio/1.2")
	for _, tc := range []struct {
		name string
		ctx  context.Context
		want string
	}{
		{name: "default", ctx: parent, want: "radiobrowser-go/0.2.0"},
		{name: "application", ctx: app, want: "MyRadio/1.2"},
		{name: "another application", ctx: rb.WithUserAgent(parent, "OtherRadio/2.0"), want: "OtherRadio/2.0"},
		{name: "derived context", ctx: context.WithValue(app, struct{}{}, true), want: "MyRadio/1.2"},
		{name: "nested override", ctx: rb.WithUserAgent(app, "OtherRadio/2.0"), want: "OtherRadio/2.0"},
		{name: "empty resets default", ctx: rb.WithUserAgent(app, ""), want: rb.DefaultUserAgent},
		{name: "whitespace resets default", ctx: rb.WithUserAgent(app, " \t "), want: rb.DefaultUserAgent},
		{name: "trim whitespace", ctx: rb.WithUserAgent(parent, " \tMyRadio/1.2 \t"), want: "MyRadio/1.2"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, method := range []string{http.MethodGet, http.MethodPost} {
				t.Run(method, func(t *testing.T) {
					received := make(chan string, 1)
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						received <- r.UserAgent()
						w.Header().Set("Content-Type", "application/json")
						if r.Method == http.MethodGet {
							_, _ = w.Write([]byte(`[]`))
						} else {
							_, _ = w.Write([]byte(`{"ok":true}`))
						}
					}))
					defer server.Close()

					var err error
					if method == http.MethodGet {
						_, err = rb.Stations(tc.ctx, server.URL, rb.StationsOptions{})
					} else {
						_, err = rb.AddStation(tc.ctx, server.URL, rb.AddStationRequest{
							Name: "My Station",
							URL:  "https://example.com/stream",
						})
					}
					if err != nil {
						t.Fatalf("request failed: %v", err)
					}
					if got := <-received; got != tc.want {
						t.Errorf("User-Agent = %q, want %q", got, tc.want)
					}
				})
			}
		})
	}
}

func TestWithUserAgentPreservesCancellation(t *testing.T) {
	deadline := time.Now().Add(time.Minute)
	parent, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	ctx := rb.WithUserAgent(parent, "MyRadio/1.2")
	if got, ok := ctx.Deadline(); !ok || !got.Equal(deadline) {
		t.Fatalf("deadline = %v, %v; want %v, true", got, ok, deadline)
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Errorf("context error = %v, want %v", ctx.Err(), context.Canceled)
	}
}
