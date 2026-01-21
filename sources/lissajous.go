package sources

import (
	"math"
)

type ScopeParams struct {
	Fx, Fy float64
	Phase  float64
	Gain   float64
	Decay  float64
}

type Lissajous struct {
	Params  *ScopeParams
	sampleT float64
	rate    float64
}

func NewLissajous(params *ScopeParams, rate float64) *Lissajous {
	return &Lissajous{
		Params: params,
		rate:   rate,
	}
}

func (l *Lissajous) Update(dt float64) {}
func (l *Lissajous) Emit(n int, out []Sample) int {
	p := l.Params
	dt := 1.0 / l.rate
	for i := range n {
		t := l.sampleT
		out[i] = Sample{
			T: t,
			XY: XY{
				X: math.Sin(2*math.Pi*p.Fx*t + p.Phase),
				Y: math.Sin(2 * math.Pi * p.Fy * t),
			},
		}
		l.sampleT += dt
	}
	return n
}
