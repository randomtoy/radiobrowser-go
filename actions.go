package radiobrowser

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/randomtoy/radiobrowser-go/internal/transport"
)

func VoteForStation(ctx context.Context, baseURL string, stationUUID string) (*VoteResponse, error) {
	stationUUID = strings.TrimSpace(stationUUID)
	if stationUUID == "" {
		return nil, fmt.Errorf("stationuuid is empty")
	}

	path := "/json/vote/" + url.PathEscape(stationUUID)

	var out VoteResponse
	if err := transport.GetJSON(ctx, baseURL, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type AddStationRequest struct {
	Name string
	URL  string

	Homepage      string
	Favicon       string
	CountryCode   string
	ISO3166_2     string
	LanguageCodes string
	// Deprecated: use ISO3166_2 instead.
	State string
	// Deprecated: use LanguageCodes instead.
	Language string
	Tags     string

	GeoLat  *float64
	GeoLong *float64
}

func (r AddStationRequest) validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is empty")
	}
	if strings.TrimSpace(r.URL) == "" {
		return fmt.Errorf("url is empty")
	}
	return validateCoordinates(r.GeoLat, r.GeoLong)
}

func (r AddStationRequest) toForm() url.Values {
	f := url.Values{}

	f.Set("name", strings.TrimSpace(r.Name))
	f.Set("url", strings.TrimSpace(r.URL))

	if s := strings.TrimSpace(r.Homepage); s != "" {
		f.Set("homepage", s)
	}
	if s := strings.TrimSpace(r.Favicon); s != "" {
		f.Set("favicon", s)
	}
	if s := strings.TrimSpace(r.CountryCode); s != "" {
		f.Set("countrycode", s)
	}
	if s := strings.TrimSpace(r.ISO3166_2); s != "" {
		f.Set("iso_3166_2", s)
	}
	if s := strings.TrimSpace(r.LanguageCodes); s != "" {
		f.Set("languagecodes", s)
	}
	if s := strings.TrimSpace(r.State); s != "" {
		f.Set("state", s)
	}
	if s := strings.TrimSpace(r.Language); s != "" {
		f.Set("language", s)
	}
	if s := strings.TrimSpace(r.Tags); s != "" {
		f.Set("tags", s)
	}

	if r.GeoLat != nil {
		f.Set("geo_lat", strconv.FormatFloat(*r.GeoLat, 'f', -1, 64))
	}
	if r.GeoLong != nil {
		f.Set("geo_long", strconv.FormatFloat(*r.GeoLong, 'f', -1, 64))
	}

	return f
}

func AddStation(ctx context.Context, baseURL string, req AddStationRequest) (*AddStationResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	var out AddStationResponse
	if err := transport.PostFormJSON(ctx, baseURL, "/json/add", req.toForm(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
