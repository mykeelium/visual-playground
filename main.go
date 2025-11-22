package main

import (
	"fmt"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/mykeelium/visual-playground/collatz"
	"github.com/mykeelium/visual-playground/primitives"
	"golang.org/x/image/colornames"
)

func main() {
	tree := collatz.BuildTree(100)

	fmt.Println("tree:")
	collatz.PrintOrganicTree(&tree)
	opengl.Run(run)
}

func drawEntities(imd *imdraw.IMDraw, entities []*primitives.Entity) {
	for _, e := range entities {
		e.Draw(imd)
	}
}

func updateEntities(entities []*primitives.Entity) {
	for _, e := range entities {
		e.Update(float64(1) / 30)
	}
}

func run() {
	var width float64 = 1024
	var height float64 = 768
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
	circles := []*primitives.Entity{}

	var startX float64 = 400
	for range 10 {
		var startY float64 = 750
		for range 2 {
			circles = append(circles, primitives.NewCircleEntity(startX, startY, 10, 0, 1, 0, 0))
		}
		startX += 40
	}

	gravity := primitives.Float2{X: 0, Y: -9.8}

	for !win.Closed() {
		imd.Clear()
		win.Clear(colornames.Black)
		primitives.ApplyGravity(circles, gravity)
		updateEntities(circles)
		drawEntities(imd, circles)
		imd.Draw(win)
		win.Update()
	}
}
