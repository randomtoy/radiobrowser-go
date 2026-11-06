package radiobrowser

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

type CountriesOptions struct {
	Order      string // "name" or "stationcount"
	Reverse    bool
	HideBroken bool
	Offset     int
}

func (o CountriesOptions) toQuery() url.Values {
	q := url.Values{}

	if o.Order != "" {
		q.Set("order", o.Order)
	}
	if o.Reverse {
		q.Set("reverse", "true")
	}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}

	return q
}

type CodecsOptions struct {
	Order      string
	Reverse    bool
	HideBroken bool
	Offset     int
}

func (o CodecsOptions) toQuery() url.Values {
	q := url.Values{}

	if o.Order != "" {
		q.Set("order", o.Order)
	}
	if o.Reverse {
		q.Set("reverse", "true")
	}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}

	return q
}

type StatesOptions struct {
	Order      string
	Reverse    bool
	HideBroken bool
	Offset     int
}

func (o StatesOptions) toQuery() url.Values {
	q := url.Values{}

	if o.Order != "" {
		q.Set("order", o.Order)
	}
	if o.Reverse {
		q.Set("reverse", "true")
	}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}

	return q
}

type LanguagesOptions struct {
	Order      string
	Reverse    bool
	HideBroken bool
	Offset     int
}

func (o LanguagesOptions) toQuery() url.Values {
	q := url.Values{}

	if o.Order != "" {
		q.Set("order", o.Order)
	}
	if o.Reverse {
		q.Set("reverse", "true")
	}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}

	return q
}

type TagsOptions struct {
	Order      string
	Reverse    bool
	HideBroken bool
	Offset     int
}

func (o TagsOptions) toQuery() url.Values {
	q := url.Values{}

	if o.Order != "" {
		q.Set("order", o.Order)
	}
	if o.Reverse {
		q.Set("reverse", "true")
	}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}

	return q
}

type StreamingServersOptions struct {
	HideBroken bool
}

func (o StreamingServersOptions) toQuery() url.Values {
	q := url.Values{}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}
	return q
}

func Countries(ctx context.Context, baseURL, filter string, opts CountriesOptions) ([]Country, error) {
	path := "/json/countries"
	filter = strings.TrimSpace(filter)
	if filter != "" {
		path += "/" + url.PathEscape(filter)
	}

	var out []Country
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func Codecs(
	ctx context.Context,
	baseURL string,
	filter string,
	opts CodecsOptions,
) ([]Codec, error) {
	path := "/json/codecs"
	filter = strings.TrimSpace(filter)
	if filter != "" {
		path += "/" + url.PathEscape(filter)
	}

	var out []Codec
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func States(
	ctx context.Context,
	baseURL string,
	country string,
	filter string,
	opts StatesOptions,
) ([]State, error) {
	path := "/json/states"

	country = strings.TrimSpace(country)
	filter = strings.TrimSpace(filter)

	if country != "" {
		path += "/" + url.PathEscape(country) + "/" + url.PathEscape(filter)
	} else if filter != "" {
		path += "/" + url.PathEscape(filter)
	}

	var out []State
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func Languages(
	ctx context.Context,
	baseURL string,
	filter string,
	opts LanguagesOptions,
) ([]Language, error) {
	path := "/json/languages"
	filter = strings.TrimSpace(filter)
	if filter != "" {
		path += "/" + url.PathEscape(filter)
	}

	var out []Language
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func Tags(
	ctx context.Context,
	baseURL string,
	filter string,
	opts TagsOptions,
) ([]Tag, error) {
	path := "/json/tags"
	filter = strings.TrimSpace(filter)
	if filter != "" {
		path += "/" + url.PathEscape(filter)
	}

	var out []Tag
	if err := transport.GetJSON(ctx, baseURL, path, opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func StreamingServers(
	ctx context.Context,
	baseURL string,
	opts StreamingServersOptions,
) ([]StreamingServer, error) {
	var out []StreamingServer
	if err := transport.GetJSON(ctx, baseURL, "/json/streamingservers", opts.toQuery(), &out); err != nil {
		return nil, err
	}
	return out, nil
}
