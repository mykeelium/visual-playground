// Package sources is used to define different sources that produce samples. Samples can be produced in different ways
package sources

type XY struct{ X, Y float64 }
type Sample struct {
	T  float64 // seconds
	XY XY
	V  float64
}

type Source interface {
	// 	Reset()
	Update(dt float64)
	Emit(n int, out []Sample) int
}
