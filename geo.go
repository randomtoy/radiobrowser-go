package radiobrowser

import (
	"fmt"
	"math"
)

func validateCoordinates(lat, long *float64) error {
	if lat != nil && (math.IsNaN(*lat) || math.IsInf(*lat, 0) || *lat < -90 || *lat > 90) {
		return fmt.Errorf("geo_lat must be a finite number between -90 and 90")
	}
	if long != nil && (math.IsNaN(*long) || math.IsInf(*long, 0) || *long < -180 || *long > 180) {
		return fmt.Errorf("geo_long must be a finite number between -180 and 180")
	}
	return nil
}

func (o StationSearchOptions) validate() error {
	if err := validateCoordinates(o.GeoLat, o.GeoLong); err != nil {
		return err
	}
	if d := o.GeoDistance; d != nil && (math.IsNaN(*d) || math.IsInf(*d, 0) || *d < 0) {
		return fmt.Errorf("geo_distance must be a finite non-negative number of meters")
	}
	return nil
}
