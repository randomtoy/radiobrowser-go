package radiobrowser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMetadataStationCounts(t *testing.T) {
	for _, model := range []struct {
		name string
		new  func() any
	}{
		{"country", func() any { return &Country{} }},
		{"codec", func() any { return &Codec{} }},
		{"state", func() any { return &State{} }},
		{"language", func() any { return &Language{} }},
		{"tag", func() any { return &Tag{} }},
	} {
		t.Run(model.name, func(t *testing.T) {
			for _, count := range []string{`12`, `"12"`, `0`, `"0"`} {
				value := model.new()
				payload := fmt.Sprintf(`{"name":"test","stationcount":%s,"country":"NL","iso_3166_1":"NL","iso_639":"eng"}`, count)
				if err := json.Unmarshal([]byte(payload), value); err != nil {
					t.Fatalf("count %s: %v", count, err)
				}
				encoded, err := json.Marshal(value)
				if err != nil {
					t.Fatal(err)
				}
				var fields map[string]any
				if err := json.Unmarshal(encoded, &fields); err != nil {
					t.Fatal(err)
				}
				want := float64(12)
				if count == `0` || count == `"0"` {
					want = 0
				}
				if fields["name"] != "test" || fields["stationcount"] != want {
					t.Fatalf("lost metadata fields: %s", encoded)
				}
			}
			for _, count := range []string{`true`, `1.5`, `"bad"`, `"1.5"`, `{}`, `9223372036854775808`} {
				if err := json.Unmarshal([]byte(`{"stationcount":`+count+`}`), model.new()); err == nil {
					t.Errorf("accepted invalid count %s", count)
				}
			}
		})
	}
}

func TestStationNullableCoordinatesAndNewFields(t *testing.T) {
	var station Station
	if err := json.Unmarshal([]byte(`{"geo_lat":null,"geo_long":null,"geo_distance":null}`), &station); err != nil {
		t.Fatal(err)
	}
	if station.GeoLat != nil || station.GeoLong != nil || station.GeoDistance != nil {
		t.Fatal("null coordinates must remain unknown")
	}
	if err := json.Unmarshal([]byte(`{"geo_lat":0,"geo_long":0,"geo_distance":0,"iso_3166_2":"US-NY","ssl_error":1,"lastlocalchecktime_iso8601":"2026-11-06T00:00:00Z"}`), &station); err != nil {
		t.Fatal(err)
	}
	if station.GeoLat == nil || *station.GeoLat != 0 || station.GeoLong == nil || *station.GeoLong != 0 || station.GeoDistance == nil || *station.GeoDistance != 0 {
		t.Fatal("zero coordinates and distance must remain present")
	}
	if station.ISO3166_2 != "US-NY" || station.SSLError != 1 || station.LastLocalCheckTimeISO != "2026-11-06T00:00:00Z" {
		t.Fatalf("missing station fields: %+v", station)
	}
}

func TestClickCounterFullResponse(t *testing.T) {
	for _, ok := range []string{`true`, `"true"`} {
		var response StationClickCounterResponse
		data := `{"ok":` + ok + `,"message":"retrieved station url","stationuuid":"s1","name":"Station","url":"https://example.com/stream"}`
		if err := json.Unmarshal([]byte(data), &response); err != nil {
			t.Fatal(err)
		}
		if !response.OK || response.Message == "" || response.StationUUID != "s1" || response.Name != "Station" || response.URL != "https://example.com/stream" {
			t.Fatalf("incomplete click response: %+v", response)
		}
	}
}
