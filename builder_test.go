package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathParsing(t *testing.T) {
	rstStr1 := `[{"id": 2251799813685425, "label": "Person", "properties": {"name": "Smith"}}::vertex, 
	{"id": 2533274790396576, "label": "workWith", "end_id": 2251799813685425, "start_id": 2251799813685424, "properties": {"weight": 3}}::edge, 
	{"id": 2251799813685424, "label": "Person", "properties": {"name": "Joe"}}::vertex]::path`

	rstStr2 := `[{"id": 2251799813685424, "label": "Person", "properties": {"name": "Joe"}}::vertex, 
	{"id": 2533274790396576, "label": "workWith", "end_id": 2251799813685425, "start_id": 2251799813685424, "properties": {"weight": 3}}::edge, 
	{"id": 2251799813685425, "label": "Person", "properties": {"name": "Smith"}}::vertex]::path`

	rstStr3 := `[{"id": 2251799813685424, "label": "Person", "properties": {"name": "Joe"}}::vertex, 
	{"id": 2533274790396579, "label": "workWith", "end_id": 2251799813685426, "start_id": 2251799813685424, "properties": {"weight": 5}}::edge, 
	{"id": 2251799813685426, "label": "Person", "properties": {"name": "Jack"}}::vertex]::path`

	unmarshaler := NewAGUnmarshaler()
	entity1, _ := unmarshaler.unmarshal(rstStr1)
	entity2, _ := unmarshaler.unmarshal(rstStr2)
	entity3, _ := unmarshaler.unmarshal(rstStr3)

	assert.Equal(t, entity1.GType(), entity2.GType(), "Type Check")
	p1 := entity1.(*Path)
	p2 := entity2.(*Path)
	p3 := entity3.(*Path)

	assert.Equal(t, p1.end.properties["name"], p2.start.properties["name"])
	assert.Equal(t, p2.start.properties["name"], p3.start.properties["name"])

	// fmt.Println(entity1)
	// fmt.Println(entity2)
	// fmt.Println(entity3)
}

func TestVertexParsing(t *testing.T) {
	rstStr := `{"id": 2251799813685425, "label": "Person", 
		"properties": {"name": "Smith", "numInt":123, "numFloat": 384.23424, "yn":true, "nullVal": null}}::vertex`

	unmarshaler := NewAGUnmarshaler()
	entity, _ := unmarshaler.unmarshal(rstStr)

	// fmt.Println(entity)
	assert.Equal(t, G_VERTEX, entity.GType())
}

func TestNormalValueParsing(t *testing.T) {
	mapStr := `{"name": "Smith", "num":123, "yn":true}`
	arrStr := `["name", "Smith", "num", 123, "yn", true]`
	// strStr1 := `abcd`
	strStr2 := `"abcd"`
	intStr := `1234`
	floatStr := `1234.56789`
	boolStr := `true`

	unmarshaler := NewAGUnmarshaler()
	mapv, _ := unmarshaler.unmarshal(mapStr)
	arrv, _ := unmarshaler.unmarshal(arrStr)
	// str1 := unmarshaler.unmarshal(strStr1)
	str2, _ := unmarshaler.unmarshal(strStr2)
	intv, _ := unmarshaler.unmarshal(intStr)
	fl, _ := unmarshaler.unmarshal(floatStr)
	b, _ := unmarshaler.unmarshal(boolStr)

	assert.Equal(t, G_MAP, mapv.GType())
	assert.Equal(t, G_ARR, arrv.GType())
	// assert.Equal(t, G_STR, str1.GType())
	assert.Equal(t, G_STR, str2.GType())
	assert.Equal(t, G_INT, intv.GType())
	assert.Equal(t, G_FLOAT, fl.GType())
	assert.Equal(t, G_BOOL, b.GType())

	// assert.Equal(t, str1.(*SimpleEntity).Value(), str2.(*SimpleEntity).Value())

	// fmt.Println(mapv)
	// fmt.Println(arrv)
	// // fmt.Println(str1)
	// fmt.Println(str2)
	// fmt.Println(intv)
	// fmt.Println(fl)
	// fmt.Println(b)
}
