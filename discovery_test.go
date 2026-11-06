package radiobrowser

import (
	"context"
	"errors"
	"net"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestDiscoverServers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	lookup := func(gotCtx context.Context, service, proto, name string) (string, []*net.SRV, error) {
		if gotCtx != ctx || service != "api" || proto != "tcp" || name != "radio-browser.info" {
			t.Fatalf("unexpected DNS lookup: ctx=%v service=%q proto=%q name=%q", gotCtx, service, proto, name)
		}
		return "_api._tcp.radio-browser.info.", []*net.SRV{
			{Target: "DE1.api.radio-browser.info.", Port: 443},
			{Target: "de1.api.radio-browser.info", Port: 443},
			{Target: "de2.api.radio-browser.info.", Port: 443},
			{Target: "de2.api.radio-browser.info.", Port: 8443},
			{Target: "de2.api.radio-browser.info.", Port: 8443},
			nil,
			{Target: ".", Port: 443},
			{Target: "", Port: 443},
			{Target: "disabled.example.", Port: 0},
		}, nil
	}
	got, err := discoverServers(ctx, lookup)
	if err != nil {
		t.Fatal(err)
	}
	// The order is intentionally random; compare the complete set of usable URLs.
	sort.Strings(got)
	want := []string{
		"https://de1.api.radio-browser.info",
		"https://de2.api.radio-browser.info",
		"https://de2.api.radio-browser.info:8443",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("servers = %v, want %v", got, want)
	}
}

func TestDiscoverServersNoUsableRecords(t *testing.T) {
	for _, records := range [][]*net.SRV{
		nil,
		{},
		{nil, {Target: ".", Port: 443}, {Target: "", Port: 443}, {Target: "disabled.example.", Port: 0}},
	} {
		got, err := discoverServers(context.Background(), func(context.Context, string, string, string) (string, []*net.SRV, error) {
			return "", records, nil
		})
		if err == nil || got != nil {
			t.Fatalf("records %v: servers = %v, err = %v; want nil servers and an error", records, got, err)
		}
	}
}

func TestDiscoverServersResolverError(t *testing.T) {
	want := &net.DNSError{Err: "temporary failure", Name: "_api._tcp.radio-browser.info", IsTemporary: true}
	got, err := discoverServers(context.Background(), func(context.Context, string, string, string) (string, []*net.SRV, error) {
		return "", []*net.SRV{{Target: "partial.example.", Port: 443}}, want
	})
	var dnsErr *net.DNSError
	if got != nil || !errors.Is(err, want) || !errors.As(err, &dnsErr) || !dnsErr.Temporary() {
		t.Fatalf("servers = %v, err = %v; want nil servers and wrapped DNS error", got, err)
	}
}

func TestDiscoverServersInvalidContext(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	expired, stop := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer stop()
	for _, tc := range []struct {
		name string
		ctx  context.Context
		want error
	}{
		{"nil", nil, nil},
		{"canceled", canceled, context.Canceled},
		{"expired", expired, context.DeadlineExceeded},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DiscoverServers(tc.ctx)
			if err == nil || got != nil || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("servers = %v, err = %v; want nil servers and error %v", got, err, tc.want)
			}
		})
	}
}

func TestDiscoverServersCancellationDuringLookup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	got, err := discoverServers(ctx, func(gotCtx context.Context, _, _, _ string) (string, []*net.SRV, error) {
		cancel()
		<-gotCtx.Done()
		return "", nil, gotCtx.Err()
	})
	if got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("servers = %v, err = %v; want nil servers and context.Canceled", got, err)
	}
}
