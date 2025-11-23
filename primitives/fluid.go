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
