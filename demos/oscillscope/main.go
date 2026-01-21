package main

import (
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mykeelium/visual-playground/engines"
	"github.com/mykeelium/visual-playground/meshes"
	"github.com/mykeelium/visual-playground/renderers"
	"github.com/mykeelium/visual-playground/sources"
	"github.com/mykeelium/visual-playground/views"
)

func handleInput(win *opengl.Window, p *sources.ScopeParams, dt float64) {
	step := 0.5 * dt

	if win.Pressed(pixel.KeyQ) {
		p.Fx += step
	}
	if win.Pressed(pixel.KeyA) {
		p.Fx -= step
	}

	if win.Pressed(pixel.KeyW) {
		p.Fy += step
	}
	if win.Pressed(pixel.KeyS) {
		p.Fy -= step
	}

	if win.Pressed(pixel.KeyE) {
		p.Phase += step
	}
	if win.Pressed(pixel.KeyD) {
		p.Phase -= step
	}

	if win.Pressed(pixel.KeyR) {
		p.Gain += step
	}
	if win.Pressed(pixel.KeyF) {
		p.Gain -= step
	}

	if win.Pressed(pixel.KeyT) {
		p.Decay += step
	}
	if win.Pressed(pixel.KeyG) {
		p.Decay -= step
	}
}

func runOld() {
	rate := 12000
	screen := views.NewPixelScreen(opengl.WindowConfig{
		Title:  "Lissajous",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	})

	scopeParams := sources.ScopeParams{
		Gain:  1.0,
		Decay: 0.96,
		Fx:    3.0,
		Fy:    2.0,
		Phase: 0.0,
	}

	source := sources.NewLissajous(&scopeParams, float64(rate))

	renderer := renderers.OscilloscopeRenderer{
		Params: &scopeParams,
	}

	engine := engines.New(
		source,
		engines.WithSampleRate(float64(rate)),
	)

	for !screen.Window().Closed() {
		dt := screen.DT()
		handleInput(screen.Window(), &scopeParams, dt)

		engine.Step(dt)

		renderer.BeginFrame(dt)
		renderer.Draw(engine.Samples())
		renderer.EndFrame(screen.Window())

		screen.Present()
	}
}

func run() {
	rate := 12000

	screen := views.NewPixelScreen(opengl.WindowConfig{
		Title:  "Lissajous",
		Bounds: pixel.R(0, 0, 1920, 1080),
		VSync:  true,
	})
	canvas := opengl.NewCanvas(screen.Window().Bounds())
	screen.SetCanvas(canvas)

	scopeParams := sources.ScopeParams{
		Gain:  1.0,
		Decay: 0.96,
		Fx:    3.0,
		Fy:    2.0,
		Phase: 0.0,
	}

	source := sources.NewLissajous(&scopeParams, float64(rate))

	engine := engines.New(
		source,
		engines.WithSampleRate(float64(rate)),
	)

	// --- rendering setup ---
	tileW := 240.0
	tileH := 135.0
	cols := 8
	rows := 8

	meshRegistry := meshes.NewMeshRegistry()
	oscMeshID := meshRegistry.Register(meshes.Mesh{Vertices: nil, Mode: meshes.DrawModeLine})

	scope :=
		renderers.Tile(
			renderers.Oscilloscope(oscMeshID),
			tileW, tileH,
			cols, rows,
		)

	renderer := &renderers.GraphRenderer{
		Root:    scope,
		Backend: renderers.NewIMDrawBackend(meshRegistry),
	}

	start := time.Now()

	for !screen.Window().Closed() {
		dt := screen.DT()
		handleInput(screen.Window(), &scopeParams, dt)

		engine.Step(dt)
		samples := engine.Samples()

		mesh := meshes.BuildOscilloscopeMesh(samples, &scopeParams, tileW, tileH)
		meshRegistry.Update(oscMeshID, mesh)

		fc := &renderers.FrameContext{
			Target:  screen.Window(),
			Time:    time.Since(start).Seconds(),
			Delta:   dt,
			Size:    screen.Window().Bounds().Size(),
			Samples: samples,
		}

		renderer.Render(fc)

		screen.Present()
	}
}

func main() {
	opengl.Run(run)
}
