package main

import (
	"math"
	"math/rand"
	"sync"
)

type Boid struct {
	position Vect
	velocity Vect
}

func NewRandom(maxWidth, maxHeight int) *Boid {
	return &Boid{
		Vect{
			rand.Float64() * float64(maxWidth),
			rand.Float64() * float64(maxHeight),
		},
		Vect{
			rand.Float64() * 0.01,
			rand.Float64() * 0.01,
		},
	}
}

func NewRandomGroup(count, maxWidth, maxHeight int) []*Boid {
	boidGroup := make([]*Boid, 0)

	for {
		if count < 0 {
			break
		}
		boidGroup = append(
			boidGroup,
			&Boid{
				Vect{
					rand.Float64() * float64(maxWidth),
					rand.Float64() * float64(maxHeight),
				},
				Vect{1, 2},
			},
		)
		count -= 1
	}

	return boidGroup
}

func (p *Boid) update() {
	p.position.x += p.velocity.x
	p.position.y += p.velocity.y
}

func (p *Boid) RoundedUpdate(w, h int) {
	if p.position.x+p.velocity.x < 0 {
		p.position.x = float64(w) - math.Mod((p.position.x+p.velocity.x), float64(w))
	} else {
		p.position.x = math.Mod((p.position.x + p.velocity.x), float64(w))
	}

	if p.position.y+p.velocity.y < 0 {
		p.position.y = float64(h) - math.Mod((p.position.y+p.velocity.y), float64(h))
	} else {
		p.position.y = math.Mod((p.position.y + p.velocity.y), float64(h))
	}
}

type BoidWorker struct {
	boid    *Boid
	config  *Config
	desire  *Vect
	friends map[int]Boid
	out     chan Vect
	id      int
}

func NewWorker(w, h int, friends map[int]Boid, id int, c *Config, d *Vect) *BoidWorker {
	return &BoidWorker{
		NewRandom(w, h),
		c,
		d,
		friends,
		make(chan Vect, 1),
		id,
	}
}

func (bw *BoidWorker) Run(l *sync.RWMutex) {
	for {
		l.RLock()
		friends := FriendsWithoutMe(bw.friends, bw.id)
		l.RUnlock()

		r1 := Rule1(friends, bw.boid, bw.config.rule1)
		r2 := Rule2(friends, bw.boid, bw.config.rule2)
		r3 := Rule3(friends, bw.boid, bw.config.rule3)
		r4 := Rule4(bw.desire, bw.boid)

		bw.boid.velocity = bw.boid.velocity.AddAll(r4, r3, r2, r1).Clamp()
		bw.boid.RoundedUpdate(640, 480)

		bw.out <- bw.boid.position
	}
}

func FriendsWithoutMe(all map[int]Boid, id int) []Boid {
	out := make([]Boid, 0)
	for i, v := range all {
		if id != i {
			out = append(out, v)
		}
	}
	return out
}

func mapKeys[T comparable, V comparable](m map[T]V) []V {
	out := make([]V, 0)

	for _, v := range m {
		out = append(out, v)
	}

	return out
}
