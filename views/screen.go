// Package views handles some fo the logic pertaining to screens, as well as controls
package views

import (
	"image/color"
	"time"

	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mykeelium/visual-playground/sources"
)

type Screen interface {
	Window() *opengl.Window
	DT() float64
	Clear()
	Present()
}

type Plotter interface {
	PlotXY(screen Screen, pts []sources.XY, intensity float64)
	PlotTime(screen Screen, vs []float64, intensity float64)
}

type PixelScreen struct {
	win    *opengl.Window
	lastT  time.Time
	dt     float64
	canvas *opengl.Canvas
}

func NewPixelScreen(cfg opengl.WindowConfig) *PixelScreen {
	win, _ := opengl.NewWindow(cfg)
	return &PixelScreen{
		win:   win,
		lastT: time.Now(),
	}
}

func (s *PixelScreen) SetCanvas(canvas *opengl.Canvas) {
	s.canvas = canvas
}

func (s *PixelScreen) Window() *opengl.Window { return s.win }

func (s *PixelScreen) DT() float64 {
	now := time.Now()
	s.dt = now.Sub(s.lastT).Seconds()
	s.lastT = now
	return s.dt
}

func (s *PixelScreen) Clear()   { s.win.Clear(color.Black) }
func (s *PixelScreen) Present() { s.win.Update() }
