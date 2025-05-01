package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type World struct {
	workers []*BoidWorker
	current map[int]Boid
	frame   []Vect
	width   int
	height  int
	config  *Config
	lock    sync.RWMutex
}

type Config struct {
	rule1 float64
	rule2 float64
	rule3 float64
}

// NewWorld creates a new world.
func NewWorld(width, height int, boidCount int) *World {
	dfl := &Config{
		rule1: 10.0,
		rule2: 10.0,
		rule3: 0.9,
	}
	w := &World{
		workers: make([]*BoidWorker, 0),
		current: make(map[int]Boid),
		frame:   make([]Vect, 0),
		width:   width,
		height:  height,
		config:  dfl,
	}
	for i := range boidCount {
		worker := NewWorker(width, height, w.current, i, dfl)
		w.workers = append(w.workers, worker)
		w.current[i] = *worker.boid
	}
	w.init()
	return w
}

func (w *World) init() {
	for _, worker := range w.workers {
		go worker.Run(&w.lock)
	}
}

func (w *World) Update(tick int) {
	frame := make([]Vect, 0)
	for _, worker := range w.workers {
		v := <-worker.out
		frame = append(frame, v)
	}

	w.lock.Lock()
	for _, worker := range w.workers {
		w.current[worker.id] = *worker.boid
	}
	w.lock.Unlock()
	w.frame = frame
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	clear(pix)

	for _, v := range w.frame {
		row := int(v.y)
		col := int(v.x)

		if !(row < 0 || col < 0 || row >= w.height || col >= w.width) {
			base := (row * 4 * w.width) + (col * 4)
			pix[base] = 0xFF
			pix[base+1] = 0xFF
			pix[base+2] = 0xFF
			pix[base+3] = 0xFF
		}
	}
}

const (
	screenWidth  = 2 * 320
	screenHeight = 2 * 240
)

type Game struct {
	world  *World
	pixels []byte
	tick   int
}

func (g *Game) Update() error {
	switch {
	case ebiten.IsKeyPressed(ebiten.Key1):
		g.world.config.rule1 += 1.0
	case ebiten.IsKeyPressed(ebiten.Key2):
		g.world.config.rule1 -= 1.0
	case ebiten.IsKeyPressed(ebiten.Key3):
		g.world.config.rule2 += 0.01
	case ebiten.IsKeyPressed(ebiten.Key4):
		g.world.config.rule2 -= 0.01
	case ebiten.IsKeyPressed(ebiten.Key5):
		g.world.config.rule3 += 0.01
	case ebiten.IsKeyPressed(ebiten.Key6):
		g.world.config.rule3 -= 0.01
	}

	if g.tick%1 == 0 {
		g.world.Update(g.tick)
	}

	g.tick += 1
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.world.Draw(g.pixels)
	screen.WritePixels(g.pixels)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("R1: %f R2: %f R3: %f", g.world.config.rule1, g.world.config.rule2, g.world.config.rule3))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()), screenWidth-90, screenHeight-40)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		world: NewWorld(screenWidth, screenHeight, 500),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of Life (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
