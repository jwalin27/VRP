package models

import (
	"math"
)

const DriverMaxTime = 12.0 * 60
const DriverCost = 500.0

type Point struct {
	X float64
	Y float64
}

func StartPoint() Point {
	return Point{X: 0, Y: 0}
}

func (p Point) DistanceTo(other Point) float64 {
	return math.Hypot(other.X-p.X, other.Y-p.Y)
}

type Load struct {
	Number  int
	Pickup  Point
	Dropoff Point
}

// Cost returns the cost for the travel from Pickup to Dropoff for this load
func (l Load) Cost() float64 {
	return l.Pickup.DistanceTo(l.Dropoff)
}
