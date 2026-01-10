package primitives

import (
	"math"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
)

var (
	ZeroFloat2 = Float2{0, 0}
	ZeroFloat3 = Float3{0, 0, 0}
)

type Color struct {
	Red   float64
	Green float64
	Blue  float64
}

type Float3 struct {
	X float64
	Y float64
	Z float64
}

func (v Float3) Add(o Float3) Float3 {
	return Float3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Float3) Scale(f float64) Float3 {
	return Float3{v.X * f, v.Y * f, v.Z * f}
}

// Float2 is Used to handle position, motion, and force in two dimensions
type Float2 struct {
	X float64
	Y float64
}

func (v Float2) Add(o Float2) Float2 {
	return Float2{v.X + o.X, v.Y + o.Y}
}

func (v Float2) Sub(o Float2) Float2 {
	return Float2{v.X - o.X, v.Y - o.Y}
}

func (v Float2) Dot(o Float2) float64 {
	return v.X*o.X + v.Y*o.Y
}

func (v Float2) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Float2) Scale(f float64) Float2 {
	return Float2{v.X * f, v.Y * f}
}

type SPHParams struct {
	RestDensity     float64
	Stiffness       float64
	Viscosity       float64
	SmoothingRadius float64
	ParticleMass    float64
}

func DefaultSPHParams() SPHParams {
	return SPHParams{
		RestDensity:     0.02,  // Empirically calibrated for 2D with h=50, particle spacing ~10
		Stiffness:       50.0,  // Lower stiffness for stability
		Viscosity:       10.0,
		SmoothingRadius: 50.0,
		ParticleMass:    1.0,
	}
}

// SPH Style Physics - 2D Kernels
func kernelPoly6(r, h float64) float64 {
	if r < 0 || r > h {
		return 0
	}
	// 2D Poly6 normalization: 4/(π*h^8)
	coeff := 4.0 / (math.Pi * math.Pow(h, 8))
	x := h*h - r*r
	return coeff * x * x * x
}

func kernelSpikyGrad(r, h float64) float64 {
	if r < 0 || r > h {
		return 0
	}
	// 2D Spiky gradient normalization: 30/(π*h^5)
	coeff := 30.0 / (math.Pi * math.Pow(h, 5))
	x := h - r
	return coeff * x * x
}

func kernelViscosityLaplacian(r, h float64) float64 {
	if r < 0 || r > h {
		return 0
	}
	// 2D Viscosity laplacian normalization: 20/(3*π*h^5)
	coeff := 20.0 / (3.0 * math.Pi * math.Pow(h, 5))
	return coeff * (h - r)
}

func SimulateSPH(
	entities []*Entity,
	grid *SpatialGrid,
	params SPHParams,
	dt float64,
	externalAcceleration Float2,
) {
	n := len(entities)
	if n == 0 {
		return
	}

	h := params.SmoothingRadius
	m := params.ParticleMass

	grid.Clear()
	for _, e := range entities {
		grid.Insert(e)
	}

	indexOf := make(map[*Entity]int, n)
	for i, e := range entities {
		indexOf[e] = i
	}

	densities := make([]float64, n)

	for i, e := range entities {
		positionI := e.Physics.Position
		rho := 0.0

		forEachNeighbor(grid, positionI, h, func(nei *Entity) {
			_, ok := indexOf[nei]
			if !ok {
				return
			}
			positionJ := nei.Physics.Position
			diff := positionJ.Sub(positionI)
			r2 := diff.X*diff.X + diff.Y*diff.Y
			if r2 > h*h {
				return
			}
			r := math.Sqrt(r2)
			w := kernelPoly6(r, h)
			rho += m * w
		})

		// Include self-contribution to density
		rho += m * kernelPoly6(0, h)

		if rho < 1e-6 {
			rho = params.RestDensity
		}

		densities[i] = rho
	}

	pressures := make([]float64, n)
	for i := range entities {
		rho := densities[i]
		pressures[i] = params.Stiffness * (rho - params.RestDensity)
	}

	for i, e := range entities {
		positionI := e.Physics.Position
		velocityI := e.Physics.Velocity
		// rhoI := densities[i]
		pressureI := pressures[i]

		var fPressure Float2
		var fViscosity Float2

		forEachNeighbor(grid, positionI, h, func(nei *Entity) {
			if nei == e {
				return
			}
			j, ok := indexOf[nei]
			if !ok {
				return
			}

			positionJ := nei.Physics.Position
			velocityJ := nei.Physics.Velocity
			rhoJ := densities[j]
			pressureJ := pressures[j]

			diff := positionJ.Sub(positionI)
			r2 := diff.X*diff.X + diff.Y*diff.Y
			if r2 > h*h || r2 < 1e-10 {
				return
			}

			r := math.Sqrt(r2)
			direction := diff.Scale(1.0 / r)

			gradW := kernelSpikyGrad(r, h)
			// More stable symmetric pressure formulation
			scaleP := -m * (pressureI + pressureJ) / (2.0 * rhoJ) * gradW
			fPressure = fPressure.Add(direction.Scale(scaleP))

			lapW := kernelViscosityLaplacian(r, h)
			vDiff := velocityJ.Sub(velocityI)
			scaleV := params.Viscosity * m / rhoJ * lapW
			fViscosity = fViscosity.Add(vDiff.Scale(scaleV))
		})

		F := fPressure.Add(fViscosity)
		acceleration := F.Scale(1.0 / m).Add(externalAcceleration)

		velocityI = velocityI.Add(acceleration.Scale(dt))
		positionI = positionI.Add(velocityI.Scale(dt))

		e.Physics.Velocity = velocityI
		e.Physics.Position = positionI
	}
}

// The PhysicsComponent contains information for running the physics calculations for an object
type PhysicsComponent struct {
	Position     Float2
	Velocity     Float2
	Acceleration Float2
	Mass         float64
}

func (p *PhysicsComponent) ApplyForce(force Float2) {
	if p.Mass != 0 {
		p.Acceleration = p.Acceleration.Add(force.Scale(1 / p.Mass))
	}
}

func (p *PhysicsComponent) Integrate(dt float64) {
	p.Velocity = p.Velocity.Add(p.Acceleration.Scale(dt))
	p.Position = p.Position.Add(p.Velocity.Scale(dt))
	p.Acceleration = ZeroFloat2
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// The RenderComponet holds the logic for rendering objects
type RenderComponet interface {
	Draw(imd *imdraw.IMDraw, position pixel.Vec)
}

type CircleRender struct {
	Radius    float64
	Thickness float64
	Color     Color
}

func (c *CircleRender) Draw(imd *imdraw.IMDraw, position pixel.Vec) {
	imd.Color = pixel.RGB(c.Color.Red, c.Color.Green, c.Color.Blue)
	imd.Push(position)
	imd.Circle(c.Radius, c.Thickness)
}

// The Entity is a marrying of a PhysicsComponent and RenderComponent
type Entity struct {
	Physics *PhysicsComponent
	Render  RenderComponet
}

func (entity *Entity) Draw(imd *imdraw.IMDraw) {
	entity.Render.Draw(imd, pixel.Vec{X: entity.Physics.Position.X, Y: entity.Physics.Position.Y})
}

func (entity *Entity) Update(dt float64) {
	entity.Physics.Integrate(dt)
}

func NewCircleEntity(x, y, radius, thickness float64, r, g, b float64) *Entity {
	physics := &PhysicsComponent{
		Position: Float2{X: x, Y: y},
		Mass:     1.0,
	}
	render := &CircleRender{
		Radius:    radius,
		Thickness: thickness,
		Color:     Color{Red: r, Green: g, Blue: b},
	}
	return &Entity{Physics: physics, Render: render}
}
