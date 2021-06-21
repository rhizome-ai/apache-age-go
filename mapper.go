package main

import (
	"fmt"
	"reflect"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/rhizome-ai/apache-age-go/parser"
)

type AGMapper struct {
	AGUnmarshaler
}

func NewAGMapper(typeMap map[string]reflect.Type) *AGMapper {
	vcache := make(map[int64]interface{})
	if typeMap == nil {
		typeMap = make(map[string]reflect.Type)
	}
	m := AGUnmarshaler{ageParser: parser.NewAgeParser(nil),
		visitor: &MapperVisitor{UnmarshalVisitor: UnmarshalVisitor{vcache: vcache},
			typeMap: typeMap},
		errListener: NewAGErrorListener(),
		vcache:      vcache,
	}

	agm := &AGMapper{AGUnmarshaler: m}
	agm.ageParser.AddErrorListener(agm.errListener)

	return agm
}

func (m *AGMapper) PutType(label string, tp reflect.Type) {
	m.visitor.(*MapperVisitor).PutType(label, tp)
}

type MapperVisitor struct {
	UnmarshalVisitor
	typeMap map[string]reflect.Type
}

func (v *MapperVisitor) PutType(label string, tp reflect.Type) {
	v.typeMap[label] = tp
}

func (v *MapperVisitor) VisitAgeout(ctx *parser.AgeoutContext) interface{} {
	rtn := v.VisitChildren(ctx)
	return rtn
}

func (v *MapperVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	fmt.Println("MapperVisitor VisitChildren ")
	var rtn interface{}
	for _, c := range node.GetChildren() {
		pt := c.(antlr.ParseTree)
		rtn = pt.Accept(v)
	}
	return rtn
}

func (v *MapperVisitor) VisitVertex(ctx *parser.VertexContext) interface{} {
	fmt.Println("MapperVisitor VisitVertex ")
	propCtx := ctx.Properties()
	props := propCtx.Accept(v).(map[string]interface{})
	// fmt.Println(" * VisitVertex:", props)
	vid := int64(props["id"].(int))
	vertex, ok := v.vcache[vid]

	if !ok {
		vertex, err := v.mapVertex(vid, props["label"].(string), props["properties"].(map[string]interface{}))
		if err != nil {
			panic(err)
		}
		v.vcache[vid] = vertex
	}
	return vertex
}

func (v *MapperVisitor) mapVertex(vid int64, label string, properties map[string]interface{}) (interface{}, error) {
	tp, ok := v.typeMap[label]

	fmt.Println("TYPE MAP : ", label, tp)
	if !ok {
		return NewVertex(vid, label, properties), nil
	}

	value := reflect.New(tp)

	for k, v := range properties {
		field := value.FieldByName(k)
		field.Set(reflect.ValueOf(v))
	}

	return value.Elem().Interface(), nil
}
