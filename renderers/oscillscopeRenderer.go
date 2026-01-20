package renderers

import (
	"image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/mykeelium/visual-playground/meshes"
	"github.com/mykeelium/visual-playground/sources"
)

type OscilloscopeRendererMode string

const (
	OscilloscopeTimeRenderer OscilloscopeRendererMode = "time"
	OscilloscopeXYRenderer   OscilloscopeRendererMode = "xy"
)

type OscilloscopeRenderer struct {
	Params *sources.ScopeParams

	canvas *opengl.Canvas
	imd    *imdraw.IMDraw
}

func (r *OscilloscopeRenderer) BeginFrame(dt float64) {
	if r.canvas == nil {
		return
	}
	if r.imd == nil {
		r.imd = imdraw.New(nil)
	}

	decay := pixel.Clamp(r.Params.Decay, 0, 1)

	r.canvas.SetComposeMethod(pixel.ComposeOver)
	r.canvas.Clear(color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: uint8(255 * (1 - decay)),
	})

	r.imd.Clear()
}

func (r *OscilloscopeRenderer) Draw(samples []sources.Sample) {
	if r.canvas == nil {
		return // will be created in EndFrame
	}

	bounds := r.canvas.Bounds()
	center := bounds.Center()
	scale := 0.45 * min(bounds.W(), bounds.H())

	r.imd.Color = color.RGBA{
		R: uint8(255 * pixel.Clamp(r.Params.Gain, 0, 1)),
		G: 255,
		B: 255,
		A: 255,
	}

	for _, s := range samples {
		x := center.X + s.XY.X*scale
		y := center.Y + s.XY.Y*scale
		r.imd.Push(pixel.V(x, y))
		r.imd.Circle(1.0, 0)
	}

	r.imd.Draw(r.canvas)
}

func (r *OscilloscopeRenderer) EndFrame(win *opengl.Window) {
	if r.canvas == nil {
		r.canvas = opengl.NewCanvas(win.Bounds())
		r.canvas.Clear(color.Black)
	}

	r.canvas.Draw(
		win,
		pixel.IM.Moved(win.Bounds().Center()),
	)
}

func Oscilloscope(mesh meshes.MeshID) RenderFn {
	return func(ctx *RenderContext, fc *FrameContext) {
		ctx.Backend.DrawMesh(mesh, ctx.Transform)
	}
}

// func Oscilloscope(
// 	params *sources.ScopeParams,
// 	width, height float64) RenderFn {
// 	return func(ctx *RenderContext, fc *FrameContext) {
// 		im := ctx.IM
// 		im.SetMatrix(ctx.Transform)
//
// 		centerX := width / 2
// 		centerY := height / 2
// 		scale := 0.45 * min(width, height)
//
// 		for _, s := range fc.Samples {
// 			x := centerX + s.XY.X*scale
// 			y := centerY + s.XY.Y*scale
// 			im.Push(pixel.V(x, y))
// 		}
//
// 		// for x := 0.0; x < width; x += 2 {
// 		// 	y := math.Sin(x*0.02+ctx.Time) * height / 2
// 		// 	im.Push(pixel.V(x, height/2+y))
// 		// }
//
// 		im.Line(1)
// 	}
// }
