// Package renderers handles samples that and prepares the output for display
package renderers

import (
	"github.com/gopxl/pixel/v2"
	// "github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/mykeelium/visual-playground/meshes"
	"github.com/mykeelium/visual-playground/primitives"
	"github.com/mykeelium/visual-playground/sources"
)

type RenderBackend interface {
	BeginFrame(fc *FrameContext)
	DrawMesh(meshID meshes.MeshID, tranform primitives.Matrix)
	EndFrame(fc *FrameContext)
}

type Renderer interface {
	BeginFrame(ctx *FrameContext)
	Draw(ctx *FrameContext)
	EndFrame(ctx *FrameContext)
}

type RenderFn func(ctx *RenderContext, fc *FrameContext)

type GraphRenderer struct {
	Root    RenderFn
	Backend RenderBackend
}

func (r *GraphRenderer) Render(fc *FrameContext) {
	r.Backend.BeginFrame(fc)

	ctx := RenderContext{
		Backend:   r.Backend,
		Transform: primitives.IM,
		Time:      fc.Time,
	}

	r.Root(&ctx, fc)

	r.Backend.EndFrame(fc)
}

// Pixel Compatibility
type IMDrawBackend struct {
	im       *imdraw.IMDraw
	registry *meshes.MeshRegistry
}

func NewIMDrawBackend(registry *meshes.MeshRegistry) *IMDrawBackend {
	return &IMDrawBackend{
		registry: registry,
	}
}

func (b *IMDrawBackend) BeginFrame(fc *FrameContext) {
	if b.im == nil {
		b.im = imdraw.New(nil)
	}
	b.im.Clear()
}

func (b *IMDrawBackend) DrawMesh(meshID meshes.MeshID, transform primitives.Matrix) {
	mesh := b.registry.Get(meshID)
	b.im.SetMatrix(pixel.Matrix(transform))
	for _, v := range mesh.Vertices {
		b.im.Push(pixel.Vec{X: v.X, Y: v.Y})
	}
	b.im.Line(1)
}

func (b *IMDrawBackend) EndFrame(fc *FrameContext) {
	b.im.Draw(fc.Target)
}

type FrameContext struct {
	Target  pixel.Target // window or framebuffer
	Time    float64      // global time in seconds
	Delta   float64      // frame delta
	Size    pixel.Vec    // target size in pixels
	Samples []sources.Sample
}

type RenderContext struct {
	Backend   RenderBackend
	Transform primitives.Matrix
	Time      float64
}

func Tile(
	base RenderFn,
	width, height float64,
	cols, rows int,
) RenderFn {
	return func(ctx *RenderContext, fc *FrameContext) {
		for y := range rows {
			for x := range cols {

				local := *ctx // value copy

				local.Transform =
					ctx.Transform.
						Moved(primitives.Float2{
							X: float64(x) * width,
							Y: float64(y) * height,
						})

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
