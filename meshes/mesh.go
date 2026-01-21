package meshes

import "github.com/mykeelium/visual-playground/primitives"

type MeshID uint32
type DrawMode string

var DrawModeLine DrawMode = "line"

type Mesh struct {
	Vertices []primitives.Float2
	Mode     DrawMode
}

type MeshRegistry struct {
	nextID MeshID
	meshes map[MeshID]Mesh
}

func NewMeshRegistry() *MeshRegistry {
	return &MeshRegistry{
		meshes: make(map[MeshID]Mesh),
	}
}

func (r *MeshRegistry) Register(m Mesh) MeshID {
	id := r.nextID
	r.nextID++
	r.meshes[id] = m
	return id
}

func (r *MeshRegistry) Get(id MeshID) Mesh {
	return r.meshes[id]
}

func (r *MeshRegistry) Update(id MeshID, m Mesh) {
	r.meshes[id] = m
}
