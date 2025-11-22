package primitives

import (
	"math"
)

type WorldBounds struct {
	MinX float64
	MaxX float64
	MinY float64
	MaxY float64
}

type SpatialGrid struct {
	CellSize float64
	Buckets  map[[2]int][]*Entity
}

func (grid *SpatialGrid) Clear() {
	for k := range grid.Buckets {
		delete(grid.Buckets, k)
	}
}

func (grid *SpatialGrid) Insert(entity *Entity) {
	x := int(math.Floor(entity.Physics.Position.X / grid.CellSize))
	y := int(math.Floor(entity.Physics.Position.Y / grid.CellSize))

	key := [2]int{x, y}
	grid.Buckets[key] = append(grid.Buckets[key], entity)
}

func (entity *Entity) HandleBoundaryCollisions(bounds WorldBounds, restitution float64, friction float64) {
	// only handling circles for now
	circle, ok := entity.Render.(*CircleRender)
	if !ok {
		return
	}

	position := &entity.Physics.Position
	velocity := &entity.Physics.Velocity
	radius := circle.Radius

	// Bottom Floor - friction applied to simulate gravity holding the entity to the ground while moving
	if position.Y-radius < bounds.MinY {
		position.Y = bounds.MinY + radius
		velocity.Y = -velocity.Y * restitution
		velocity.X *= friction
	}

	// Top Wall
	if position.Y+radius > bounds.MaxY {
		position.Y = bounds.MaxY - radius
		velocity.Y = -velocity.Y * restitution
	}

	// Left Wall
	if position.X-radius < bounds.MinX {
		position.X = bounds.MinX + radius
		velocity.X = -velocity.X * restitution
	}

	// Right Wall
	if position.X+radius > bounds.MaxX {
		position.X = bounds.MaxX - radius
		velocity.X = -velocity.X * restitution
	}
}

func (entity *Entity) ApplyGravity(g Float2) {
	entity.Physics.ApplyForce(g.Scale(entity.Physics.Mass))
}

func (entity *Entity) HandleObjectCollisions(grid *SpatialGrid) {
	for _, cell := range grid.Buckets {
		for i := 0; i < len(cell); i++ {
			for j := i + 1; j < len(cell); j++ {
				resolveCollisions(cell[i], cell[j])
			}
		}
	}
}

func resolveCollisions(a, b *Entity) {
	// only handling circles for now
	ca, ok1 := a.Render.(*CircleRender)
	cb, ok2 := b.Render.(*CircleRender)
	if !ok1 || !ok2 {
		return
	}

	positionA := a.Physics.Position
	positionB := b.Physics.Position

	delta := positionB.Sub(positionA)
	dist := delta.Len()
	minDist := ca.Radius + cb.Radius

	if dist < minDist && dist > 0 {
		// Push them appart
		overlap := minDist - dist
		normal := delta.Scale(1 / dist)

		a.Physics.Position = a.Physics.Position.Add(normal.Scale(-overlap / 2))
		b.Physics.Position = b.Physics.Position.Add(normal.Scale(overlap / 2))

		// elastic collisions
		relVel := b.Physics.Velocity.Sub(a.Physics.Velocity)
		sepVel := relVel.Dot(normal)

		// if not going in separate directions
		if sepVel < 0 {
			impulse := normal.Scale(-2 * sepVel / (a.Physics.Mass + b.Physics.Mass))
			a.Physics.Velocity = a.Physics.Velocity.Sub(impulse.Scale(a.Physics.Mass))
			b.Physics.Velocity = b.Physics.Velocity.Add(impulse.Scale(b.Physics.Mass))
		}
	}
}

func GroupApplyGravity(e []*Entity, g Float2) {
	for _, ent := range e {
		ent.Physics.ApplyForce(g.Scale(ent.Physics.Mass))
	}
}
