// Package renderers handles samples that and prepares the output for display
package renderers

import (
	"github.com/gopxl/pixel/v2"
	// "github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/mykeelium/visual-playground/sources"
)

// type Renderer interface {
// 	BeginFrame(dt float64)
// 	Draw(samples []sources.Sample)
// 	EndFrame(target *opengl.Window)
// }

type Renderer interface {
	BeginFrame(ctx *FrameContext)
	Draw(ctx *FrameContext)
	EndFrame(ctx *FrameContext)
}

type RenderFn func(ctx *RenderContext, fc *FrameContext)

type GraphRenderer struct {
	Root RenderFn
	IM   *imdraw.IMDraw
}

func (r *GraphRenderer) BeginFrame() {
	if r.IM == nil {
		r.IM = imdraw.New(nil)
	}

	r.IM.Clear()
}

func (r *GraphRenderer) Draw(ctx *FrameContext) {
	renderCtx := RenderContext{
		Target:    ctx.Target,
		Transform: pixel.IM,
		Time:      ctx.Time,
		IM:        r.IM,
	}
	r.Root(&renderCtx, ctx)
}

func (r *GraphRenderer) EndFrame(ctx *FrameContext) {
	r.IM.Draw(ctx.Target)
}

type FrameContext struct {
	Target  pixel.Target // window or framebuffer
	Time    float64      // global time in seconds
	Delta   float64      // frame delta
	Size    pixel.Vec    // target size in pixels
	Samples []sources.Sample
}

type RenderContext struct {
	Target    pixel.Target
	IM        *imdraw.IMDraw
	Transform pixel.Matrix
	Time      float64
}

func Tile(
	base RenderFn,
	width, height float64,
	cols, rows int,
) RenderFn {
	return func(ctx *RenderContext, fc *FrameContext) {
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {

				local := *ctx // value copy

				local.Transform =
					ctx.Transform.
						Moved(pixel.V(
							float64(x)*width,
							float64(y)*height,
						))

				base(&local, fc)
			}
		}
	}
}

// func WithTileSpace(w, h float64, r RenderFn) RenderFn {
// 	return func(ctx *RenderContext) {
// 		local := *ctx
// 		local.Transform =
// 			ctx.Transform.
// 				ScaledXY(pixel.ZV, pixel.V(1/w, 1/h))
//
// 		r(&local)
// 	}
// }
