package main

import (
	"fmt"
	"testing"
)

func TestPathParsing(t *testing.T) {
	rstStr := `[{"id": 2251799813685425, "label": "Person", "properties": {"name": "Smith"}}::vertex, 
	{"id": 2533274790396576, "label": "workWith", "end_id": 2251799813685425, "start_id": 2251799813685424, "properties": {"weight": 3}}::edge, 
	{"id": 2251799813685424, "label": "Person", "properties": {"name": "Joe"}}::vertex]::path`

	unmarshaler := NewAGUnmarshaler()
	entity := unmarshaler.unmarshalAGResult(rstStr)

	fmt.Println(entity)
}

func TestVertexParsing(t *testing.T) {
	rstStr := `{"id": 2251799813685425, "label": "Person", "properties": {"name": "Smith", "num":123, "yn":true}}::vertex`

	unmarshaler := NewAGUnmarshaler()
	entity := unmarshaler.unmarshalAGResult(rstStr)

	fmt.Println(entity)
}
