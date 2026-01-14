// Package renderers handles samples that and prepares the output for display
package renderers

import (
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mykeelium/visual-playground/sources"
)

type Renderer interface {
	BeginFrame(dt float64)
	Draw(samples []sources.Sample)
	EndFrame(target *opengl.Window)
}
