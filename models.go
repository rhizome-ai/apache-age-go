package main

type GTYPE uint8

const G_VERTEX GTYPE = 1
const G_EDEGE GTYPE = 2
const G_PATH GTYPE = 3

type Entity interface {
	GType() GTYPE
}

type Node struct {
	Entity
	id         int64
	label      string
	properties map[string]interface{}
}

func newNode(id int64, label string, props map[string]interface{}) *Node {
	return &Node{id: id, label: label, properties: props}
}

func (n *Node) Id() int64 {
	return n.id
}

func (n *Node) Label() string {
	return n.label
}

func (n *Node) Prop(key string) interface{} {
	return n.properties[key]
}

// func (n *Node) String() string {
// 	fmt.Sprintf("{id: %s, label")
// 	return n.properties[key]
// }

type Vertex struct {
	*Node
}

func NewVertex(id int64, label string, props map[string]interface{}) *Vertex {
	return &Vertex{newNode(id, label, props)}
}

func (v *Vertex) GType() GTYPE {
	return G_VERTEX
}

type Edge struct {
	*Node
	start_id int64
	end_id   int64
}

func NewEdge(id int64, label string, start int64, end int64, props map[string]interface{}) *Edge {
	return &Edge{Node: newNode(id, label, props), start_id: start, end_id: end}
}

func (e *Edge) GType() GTYPE {
	return G_EDEGE
}

func (e *Edge) StartId() int64 {
	return e.start_id
}

func (e *Edge) EndId() int64 {
	return e.end_id
}

type Path struct {
	Entity
	start *Vertex
	rel   *Edge
	end   *Vertex
}

func NewPath(start *Vertex, rel *Edge, end *Vertex) *Path {
	return &Path{start: start, rel: rel, end: end}
}

func (e *Path) GType() GTYPE {
	return G_PATH
}

func (p *Path) Start() *Vertex {
	return p.start
}

func (p *Path) Rel() *Edge {
	return p.rel
}

func (p *Path) End() *Vertex {
	return p.end
}
