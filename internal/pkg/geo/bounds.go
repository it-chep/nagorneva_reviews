// Package geo provides conservative bounds for spherical radius searches.
package geo

import "math"

// Bounds returns south, north, west, east in degrees. West > east crosses
// the antimeridian. A zero radius means the existing unrestricted search.
func Bounds(lat, lon, radius float64) (south, north, west, east float64) {
	if radius <= 0 || radius >= math.Pi*6371 {
		return -90, 90, -180, 180
	}
	angular := radius / 6371
	delta := angular * 180 / math.Pi
	south, north = math.Max(-90, lat-delta), math.Min(90, lat+delta)
	if south <= -90 || north >= 90 {
		return south, north, -180, 180
	}
	lonDelta := math.Asin(math.Min(1, math.Sin(angular)/math.Cos(lat*math.Pi/180))) * 180 / math.Pi
	wrap := func(v float64) float64 { return math.Mod(v+540, 360) - 180 }
	return south, north, wrap(lon - lonDelta), wrap(lon + lonDelta)
}
