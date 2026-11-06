package radiobrowser

import (
	"encoding/json"
	"strconv"
)

// UnmarshalJSON also accepts the quoted ok flag shown in the API's click example.
func (v *StationClickCounterResponse) UnmarshalJSON(data []byte) error {
	type plain StationClickCounterResponse
	aux := struct {
		*plain
		OK json.RawMessage `json:"ok"`
	}{plain: (*plain)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.OK) == 0 {
		return nil
	}
	if aux.OK[0] == '"' {
		var text string
		if err := json.Unmarshal(aux.OK, &text); err != nil {
			return err
		}
		value, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		v.OK = value
		return nil
	}
	return json.Unmarshal(aux.OK, &v.OK)
}

// jsonInt accepts the numeric values sent by current servers and the quoted
// integers shown by older API documentation, without changing public int fields.
type jsonInt int

func (n *jsonInt) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var text string
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		value, err := strconv.Atoi(text)
		if err != nil {
			return err
		}
		*n = jsonInt(value)
		return nil
	}
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*n = jsonInt(value)
	return nil
}

// UnmarshalJSON accepts current numeric and legacy quoted station counts.
func (v *Country) UnmarshalJSON(data []byte) error {
	type plain Country
	aux := struct {
		*plain
		StationCount jsonInt `json:"stationcount"`
	}{plain: (*plain)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.StationCount = int(aux.StationCount)
	return nil
}

// UnmarshalJSON accepts current numeric and legacy quoted station counts.
func (v *Codec) UnmarshalJSON(data []byte) error {
	type plain Codec
	aux := struct {
		*plain
		StationCount jsonInt `json:"stationcount"`
	}{plain: (*plain)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.StationCount = int(aux.StationCount)
	return nil
}

// UnmarshalJSON accepts current numeric and legacy quoted station counts.
func (v *State) UnmarshalJSON(data []byte) error {
	type plain State
	aux := struct {
		*plain
		StationCount jsonInt `json:"stationcount"`
	}{plain: (*plain)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.StationCount = int(aux.StationCount)
	return nil
}

// UnmarshalJSON accepts current numeric and legacy quoted station counts.
func (v *Language) UnmarshalJSON(data []byte) error {
	type plain Language
	aux := struct {
		*plain
		StationCount jsonInt `json:"stationcount"`
	}{plain: (*plain)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.StationCount = int(aux.StationCount)
	return nil
}

// UnmarshalJSON accepts current numeric and legacy quoted station counts.
func (v *Tag) UnmarshalJSON(data []byte) error {
	type plain Tag
	aux := struct {
		*plain
		StationCount jsonInt `json:"stationcount"`
	}{plain: (*plain)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.StationCount = int(aux.StationCount)
	return nil
}
