package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mykeelium/visual-playground/engines"
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

func run() {
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

func main() {
	opengl.Run(run)
}
