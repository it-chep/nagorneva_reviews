package geo

import (
	"math"
	"testing"
)

func TestBoundsContainsCircle(t *testing.T) {
	for _, c := range [][3]float64{{60, 60, 100}, {0, 179, 500}, {70, -179, 1500}, {89, 0, 500}, {-80, 20, 2000}, {0, 0, 15000}} {
		south, north, west, east := Bounds(c[0], c[1], c[2])
		lat, lon, d := c[0]*math.Pi/180, c[1]*math.Pi/180, c[2]/6371
		for degrees := 0; degrees < 360; degrees++ {
			b := float64(degrees) * math.Pi / 180
			y := math.Asin(math.Sin(lat)*math.Cos(d) + math.Cos(lat)*math.Sin(d)*math.Cos(b))
			x := lon + math.Atan2(math.Sin(b)*math.Sin(d)*math.Cos(lat), math.Cos(d)-math.Sin(lat)*math.Sin(y))
			yd, xd := y*180/math.Pi, math.Mod(x*180/math.Pi+540, 360)-180
			insideLon := xd >= west-1e-9 && xd <= east+1e-9
			if west > east {
				insideLon = xd >= west-1e-9 || xd <= east+1e-9
			}
			if yd < south-1e-9 || yd > north+1e-9 || !insideLon {
				t.Fatalf("circle %v excludes (%f,%f) bounds %v", c, yd, xd, []float64{south, north, west, east})
			}
		}
	}
}
func TestUnrestrictedBounds(t *testing.T) {
	s, n, w, e := Bounds(60, 60, 0)
	if s != -90 || n != 90 || w != -180 || e != 180 {
		t.Fatal("zero radius must cover the world")
	}
}
