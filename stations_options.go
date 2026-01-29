package radiobrowser

import (
	"net/url"
	"strconv"
	"strings"
)

type StationOrder string

const (
	StationOrderName          StationOrder = "name"
	StationOrderURL           StationOrder = "url"
	StationOrderHomepage      StationOrder = "homepage"
	StationOrderFavicon       StationOrder = "favicon"
	StationOrderTags          StationOrder = "tags"
	StationOrderCountry       StationOrder = "country"
	StationOrderState         StationOrder = "state"
	StationOrderLanguage      StationOrder = "language"
	StationOrderVotes         StationOrder = "votes"
	StationOrderCodec         StationOrder = "codec"
	StationOrderBitrate       StationOrder = "bitrate"
	StationOrderLastCheckOK   StationOrder = "lastcheckok"
	StationOrderLastCheckTime StationOrder = "lastchecktime"
	StationOrderClickTS       StationOrder = "clicktimestamp"
	StationOrderClickCount    StationOrder = "clickcount"
	StationOrderClickTrend    StationOrder = "clicktrend"
	StationOrderChangeTS      StationOrder = "changetimestamp"
	StationOrderRandom        StationOrder = "random"
)

type StationsOptions struct {
	Order      StationOrder
	Reverse    bool
	Offset     int
	Limit      int
	HideBroken bool
}

func (o StationsOptions) toQuery() url.Values {
	q := url.Values{}

	if o.Order != "" {
		q.Set("order", string(o.Order))
	}
	if o.Reverse {
		q.Set("reverse", "true")
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}
	if o.Limit > 0 {
		q.Set("limit", strconv.Itoa(o.Limit))
	}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}

	return q
}

type StationSearchKind string

const (
	StationByUUID             StationSearchKind = "byuuid"
	StationByName             StationSearchKind = "byname"
	StationByNameExact        StationSearchKind = "bynameexact"
	StationByCodec            StationSearchKind = "bycodec"
	StationByCodecExact       StationSearchKind = "bycodecexact"
	StationByCountry          StationSearchKind = "bycountry"
	StationByCountryExact     StationSearchKind = "bycountryexact"
	StationByCountryCodeExact StationSearchKind = "bycountrycodeexact"
	StationByState            StationSearchKind = "bystate"
	StationByStateExact       StationSearchKind = "bystateexact"
	StationByLanguage         StationSearchKind = "bylanguage"
	StationByLanguageExact    StationSearchKind = "bylanguageexact"
	StationByTag              StationSearchKind = "bytag"
	StationByTagExact         StationSearchKind = "bytagexact"
)

type StationsChangedOptions struct {
	LastChangeUUID string
	Limit          int
}

func (o StationsChangedOptions) toQuery() url.Values {
	q := url.Values{}

	if s := strings.TrimSpace(o.LastChangeUUID); s != "" {
		q.Set("lastchangeuuid", s)
	}
	if o.Limit > 0 {
		q.Set("limit", strconv.Itoa(o.Limit))
	}

	return q
}

type StationSearchOptions struct {
	Name          string
	NameExact     bool
	Country       string
	CountryExact  bool
	CountryCode   string
	State         string
	StateExact    bool
	Language      string
	LanguageExact bool
	Tag           string
	TagExact      bool
	TagList       []string
	Codec         string

	BitrateMin int
	BitrateMax int

	HasGeoInfo      *bool
	HasExtendedInfo *bool
	IsHTTPS         *bool

	Order      StationOrder
	Reverse    bool
	Offset     int
	Limit      int
	HideBroken bool
}

func (o StationSearchOptions) toQuery() url.Values {
	q := url.Values{}

	if s := strings.TrimSpace(o.Name); s != "" {
		q.Set("name", s)
	}
	if o.NameExact {
		q.Set("nameExact", "true")
	}

	if s := strings.TrimSpace(o.Country); s != "" {
		q.Set("country", s)
	}
	if o.CountryExact {
		q.Set("countryExact", "true")
	}

	if s := strings.TrimSpace(o.CountryCode); s != "" {
		q.Set("countrycode", s)
	}

	if s := strings.TrimSpace(o.State); s != "" {
		q.Set("state", s)
	}
	if o.StateExact {
		q.Set("stateExact", "true")
	}

	if s := strings.TrimSpace(o.Language); s != "" {
		q.Set("language", s)
	}
	if o.LanguageExact {
		q.Set("languageExact", "true")
	}

	if s := strings.TrimSpace(o.Tag); s != "" {
		q.Set("tag", s)
	}
	if o.TagExact {
		q.Set("tagExact", "true")
	}

	if len(o.TagList) > 0 {
		var tags []string
		for _, t := range o.TagList {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			tags = append(tags, t)
		}
		if len(tags) > 0 {
			q.Set("tagList", strings.Join(tags, ","))
		}
	}

	if s := strings.TrimSpace(o.Codec); s != "" {
		q.Set("codec", s)
	}

	if o.BitrateMin > 0 {
		q.Set("bitrateMin", strconv.Itoa(o.BitrateMin))
	}
	if o.BitrateMax > 0 {
		q.Set("bitrateMax", strconv.Itoa(o.BitrateMax))
	}

	if o.HasGeoInfo != nil {
		q.Set("has_geo_info", strconv.FormatBool(*o.HasGeoInfo))
	}
	if o.HasExtendedInfo != nil {
		q.Set("has_extended_info", strconv.FormatBool(*o.HasExtendedInfo))
	}
	if o.IsHTTPS != nil {
		q.Set("is_https", strconv.FormatBool(*o.IsHTTPS))
	}

	if o.Order != "" {
		q.Set("order", string(o.Order))
	}
	if o.Reverse {
		q.Set("reverse", "true")
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}
	if o.Limit > 0 {
		q.Set("limit", strconv.Itoa(o.Limit))
	}
	if o.HideBroken {
		q.Set("hidebroken", "true")
	}

	return q
}
