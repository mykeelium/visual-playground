package main

import (
	// "fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"

	// "github.com/mykeelium/visual-playground/collatz"
	"github.com/mykeelium/visual-playground/primitives"
	"golang.org/x/image/colornames"
)

var (
	gravity     = primitives.Float2{X: 0, Y: -500}
	worldBounds = primitives.WorldBounds{
		MinY: floor,
		MinX: 0,
		MaxY: height,
		MaxX: width,
	}
	circles []*primitives.Entity
)

const (
	width        float64 = 2048
	height       float64 = 1024
	floor        float64 = 0
	resititution float64 = 0.6
	friction     float64 = 0.95
)

func CreateChaoticCircles(
	count int,
	radius float64,
	windowWidth float64,
	windowHeight float64,
	speedMin float64,
	speedMax float64,
) []*primitives.Entity {
	circles = []*primitives.Entity{}

	for range count {
		x := rand.Float64()*(windowWidth-2*radius) + radius
		y := rand.Float64()*(windowHeight-2*radius) + radius

		e := primitives.NewCircleEntity(x, y, radius, 0,
			rand.Float64(), // red
			rand.Float64(), // green
			rand.Float64(), // blue
		)

		speed := speedMin + rand.Float64()*(speedMax-speedMin)
		angle := rand.Float64() * 2 * math.Pi
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		e.Physics.Velocity = primitives.Float2(pixel.V(vx, vy))

		circles = append(circles, e)
	}

	return circles
}

func ResetSimulation() {
	circles = CreateChaoticCircles(
		500, // count
		25,  // radius
		width,
		height,
		50,  // min speed
		200, // max speed
	)
}

func main() {
	// tree := collatz.BuildTree(100)
	// fmt.Println("tree:")
	// collatz.PrintOrganicTree(&tree)
	rand.Seed(42)
	opengl.Run(run)
}

func drawEntities(imd *imdraw.IMDraw, entities []*primitives.Entity) {
	for _, e := range entities {
		e.Draw(imd)
	}
}

func updateEntities(entities []*primitives.Entity, grid *primitives.SpatialGrid, dt float64) {
	iterations := 5
	for range iterations {
		primitives.SimulateSPH(entities, grid, primitives.DefaultSPHParams(), dt/float64(iterations), gravity)
	}

	for _, e := range entities {
		// e.ApplyGravity(gravity)
		// e.Update(dt / 2)
		e.HandleBoundaryCollisions(worldBounds, resititution, friction)
		e.UpdateColorBasedOnSpeed(800)

		// grid.Insert(e)
	}

	// Multiple iteration to try and handle all collisions

	// for range 4 {
	// for _, e := range entities {
	// 	e.HandleObjectCollisions(grid)
	// }
	//}
}

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	grid := &primitives.SpatialGrid{
		CellSize: 15,
		Buckets:  map[[2]int][]*primitives.Entity{},
	}
	ResetSimulation()

	last := time.Now()

	for !win.Closed() {
		if win.JustPressed(pixel.KeyR) {
			ResetSimulation()
		}

		if win.JustPressed(pixel.KeyX) {
			win.SetClosed(true)
			continue
		}

		now := time.Now()
		dt := now.Sub(last).Seconds()
		last = now

		if win.Pressed(pixel.MouseButtonLeft) {
			mousePosition := win.MousePosition()
			primitives.ApplyCircularForce(true, circles, mousePosition, 300, 5000, dt)
		}

		if win.Pressed(pixel.MouseButtonRight) {
			mousePosition := win.MousePosition()
			primitives.ApplyCircularForce(false, circles, mousePosition, 300, 5000, dt)
		}

		// Update Simulation
		updateEntities(circles, grid, dt)

		// Draw
		imd.Clear()
		win.Clear(colornames.Black)
		drawEntities(imd, circles)
		imd.Draw(win)
		win.Update()
	}
}
