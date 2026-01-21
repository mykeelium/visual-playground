// Package meshes is used to contain specific structs and their requried values for containing meshes
package meshes

import (
	"github.com/mykeelium/visual-playground/primitives"
	"github.com/mykeelium/visual-playground/sources"
)

type OscilloscopeMesh struct {
	points []primitives.Float2
}

func BuildOscilloscopeMesh(
	samples []sources.Sample,
	params *sources.ScopeParams,
	width, height float64,
) Mesh {
	pts := make([]primitives.Float2, 0, len(samples))
	cx, cy := width/2, height/2

	for _, s := range samples {
		x := s.XY.X * params.Gain * cx
		y := s.XY.Y * params.Gain * cy
		pts = append(pts, primitives.Float2{X: cx + x, Y: cy + y})
	}

	return Mesh{Vertices: pts, Mode: DrawModeLine}
}

// func DrawMesh(
// 	mesh OscilloscopeMesh,
// 	ctx *renderers.RenderContext,
// ) {
// 	im := ctx.IM
// 	im.SetMatrix(ctx.Transform)
//
// 	for _, p := range mesh.points {
// 		im.Push(p)
// 	}
// }
