package main

import (
	"fmt"
	"reflect"
)

type GTYPE uint8

const (
	G_OTHER GTYPE = 1 + iota
	G_VERTEX
	G_EDEGE
	G_PATH
	G_STR
	G_INT
	G_FLOAT
	G_BOOL
	G_MAP
	G_ARR
)

var _TpV = reflect.TypeOf(&Vertex{})
var _TpE = reflect.TypeOf(&Edge{})
var _TpP = reflect.TypeOf(&Path{})
var _TpStr = reflect.TypeOf(string(""))
var _TpInt = reflect.TypeOf(int64(0))
var _TpFloat = reflect.TypeOf(float64(0))
var _TpBool = reflect.TypeOf(bool(false))
var _TpMap = reflect.TypeOf(map[string]interface{}{})
var _TpArr = reflect.TypeOf([]interface{}{})

type Entity interface {
	GType() GTYPE
}

func IsEntity(v interface{}) bool {
	_, ok := v.(Entity)
	return ok
}

type SimpleEntity struct {
	Entity
	typ   GTYPE
	value interface{}
}

func NewSimpleEntity(value interface{}) *SimpleEntity {
	switch value.(type) {
	case string:
		return &SimpleEntity{typ: G_STR, value: value}
	case int64:
		return &SimpleEntity{typ: G_INT, value: value}
	case float64:
		return &SimpleEntity{typ: G_FLOAT, value: value}
	case bool:
		return &SimpleEntity{typ: G_BOOL, value: value}
	case map[string]interface{}:
		return &SimpleEntity{typ: G_MAP, value: value}
	case []interface{}:
		return &SimpleEntity{typ: G_ARR, value: value}
	default:
		return &SimpleEntity{typ: G_OTHER, value: value}
	}
}

func (e *SimpleEntity) GType() GTYPE {
	return e.typ
}

func (e *SimpleEntity) Value() interface{} {
	return e.value
}

func (e *SimpleEntity) String() string {
	return fmt.Sprintf("%v", e.value)
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

type Vertex struct {
	*Node
}

func NewVertex(id int64, label string, props map[string]interface{}) *Vertex {
	return &Vertex{newNode(id, label, props)}
}

func (v *Vertex) GType() GTYPE {
	return G_VERTEX
}

func (v *Vertex) String() string {
	return fmt.Sprintf("V{id:%d, label:%s, properties:%v}", v.id, v.label, v.properties)
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

func (e *Edge) String() string {
	return fmt.Sprintf("E{id:%d, label:%s, start:%d, end:%d, properties:%v}",
		e.id, e.label, e.start_id, e.end_id, e.properties)
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

func (p *Path) String() string {
	return fmt.Sprintf("P[%v, %v, %v]",
		p.start, p.rel, p.end)
}
