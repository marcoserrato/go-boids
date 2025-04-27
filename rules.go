package main

import (
	"math"
)

func Distance(me, other Boid) float64 {
	return math.Sqrt(
		math.Pow(other.position.x-me.position.x, 2) + math.Pow(other.position.y-me.position.y, 2),
	)
}

// Cohesion: move towards the center of mass. This should be
// achievable by summing up all positions and diving by the number
// of boids. Then build a vector from the boid to the center of mass
func Rule1(all []Boid, boid *Boid, factor float64) Vect {
	centerOfMass := Vect{0, 0}
	for _, v := range all {
		centerOfMass = centerOfMass.Add(v.position)
	}

	result, _ := centerOfMass.Div(float64(len(all)))
	result = result.Min(boid.position)

	if !result.ZeroHuh() {
		result, _ = result.Div(factor)
	}

	return result
}

// Rule 2: Avoid Crowding
// - Find surrounding boids by calculating distance between.
// - Build a vector from considered boid to neighbor.
// - Subtract it from total vector.
//
// Note: considering yourself shouldn't impact the resulting vector
// since it should be 0.
func Rule2(all []Boid, boid *Boid, factor float64) Vect {
	pushAway := Vect{0, 0}
	for _, v := range all {
		distance := Distance(v, *boid)
		if distance < 20.0 {
			toNeighbor := boid.position.Min(v.position)
			toNeighbor, _ = toNeighbor.Div(distance)
			pushAway = pushAway.Add(toNeighbor)
		}
	}

	pushAway, _ = pushAway.Div(factor)

	return pushAway
}

// Rule 3: Alignment
// Find the avergage heading by summing up all velocities and
// dividing by the number of boids
func Rule3(all []Boid, boid *Boid, factor float64) Vect {
	average := Vect{0, 0}
	for _, v := range all {
		average = average.Add(v.velocity)
	}

	average, _ = average.Div(float64(len(all)))
	average, _ = average.Div(factor)

	return average
}
