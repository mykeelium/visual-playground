// Package engines is used to orastrate the combination of the source and the renderer
package engines

import (
	"github.com/mykeelium/visual-playground/renderers"
	"github.com/mykeelium/visual-playground/sources"
)

type Engine struct {
	Source      sources.Source
	Render      renderers.Renderer
	RateHz      float64
	buffer      []sources.Sample
	bufferCount int
	tAcc        float64
}

type Option func(*Engine)

func WithSampleRate(rate float64) Option {
	return func(e *Engine) {
		e.RateHz = rate
	}
}

func New(source sources.Source, opts ...Option) *Engine {
	e := &Engine{
		Source: source,
		RateHz: 120000,
		buffer: make([]sources.Sample, 0, 8192),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *Engine) Step(dt float64) {
	e.Source.Update(dt)

	// produces samples at a fixed sample rate independent of FPS
	e.tAcc += dt
	want := int(e.tAcc * e.RateHz)
	if want <= 0 {
		e.bufferCount = 0
		return
	}
	e.tAcc -= float64(want) / e.RateHz

	if cap(e.buffer) < want {
		e.buffer = make([]sources.Sample, want)
	}
	e.buffer = e.buffer[:want]
	e.bufferCount = e.Source.Emit(want, e.buffer)
}

func (e *Engine) Samples() []sources.Sample {
	return e.buffer[:e.bufferCount]
}
