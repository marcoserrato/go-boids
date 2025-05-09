package main

import (
	"fmt"
	"math"
)

type Vect struct {
	x float64
	y float64
}

func (v Vect) AddAll(others ...Vect) Vect {
	out := v
	for _, vect := range others {
		out = out.Add(vect)
	}
	return out
}

func (v Vect) MinAll(others ...Vect) Vect {
	out := v
	for _, vect := range others {
		out = out.Min(vect)
	}
	return out
}

func (v Vect) Add(other Vect) Vect {
	return Vect{v.x + other.x, v.y + other.y}
}

func (v Vect) Min(other Vect) Vect {
	return Vect{v.x - other.x, v.y - other.y}
}

func (v Vect) Avg(other Vect) Vect {
	return Vect{(v.x + other.x) / 2, (v.y + other.y) / 2}
}

func (v Vect) Div(numerator float64) (Vect, error) {
	// This doesn't work.. it's a float..
	if numerator == 0 {
		return Vect{}, fmt.Errorf("Cannot divide by 0!")
	}

	return Vect{v.x / numerator, v.y / numerator}, nil
}

func (v Vect) Mul(scalar float64) Vect {
	return Vect{v.x * scalar, v.y * scalar}
}

func (v Vect) magnitude() float64 {
	return math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
}

func (v Vect) Normalize() Vect {
	return Vect{v.x / v.magnitude(), v.y / v.magnitude()}
}

func (v Vect) ZeroHuh() bool {
	return v.x == 0 && v.y == 0
}

func VectorDistance(me, other Vect) float64 {
	return math.Sqrt(
		math.Pow(other.x-me.x, 2) + math.Pow(other.y-me.y, 2),
	)
}

func (v Vect) Clamp() Vect {
	mag := math.Sqrt(v.x*v.x + v.y*v.y)
	if mag > 3 {
		scale := 3 / mag
		return Vect{v.x * scale, v.y * scale}
	}
	return v
}
