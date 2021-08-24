package main

import (
	"database/sql"
	"fmt"

	"github.com/rhizome-ai/apache-age-go/age"
)

// Do cypher query to AGE with database/sql Tx API transaction conrol
func doWithSqlAPI(dsn string, graphName string) {

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// Confirm graph_path created
	_, err = age.GetReady(db, graphName)
	if err != nil {
		panic(err)
	}

	// Tx begin for execute create vertex
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// Create vertices with Cypher
	_, err = age.ExecCypher(tx, graphName, 0, "CREATE (n:Person {name: '%s', weight:%f})", "Joe", 67.3)
	if err != nil {
		panic(err)
	}

	_, err = age.ExecCypher(tx, graphName, 0, "CREATE (n:Person {name: '%s', weight:77.3, roles:['Dev','marketing']})", "Jack")
	if err != nil {
		panic(err)
	}

	_, err = age.ExecCypher(tx, graphName, 0, "CREATE (n:Person {name: '%s', weight:%d})", "Andy", 59)
	if err != nil {
		panic(err)
	}

	// Commit Tx
	tx.Commit()

	// Tx begin for queries
	tx, err = db.Begin()
	if err != nil {
		panic(err)
	}
	// Query cypher
	cypherCursor, err := age.ExecCypher(tx, graphName, 1, "MATCH (n:Person) RETURN n")
	if err != nil {
		panic(err)
	}
	// Unmarsal result data to Vertex row by row
	for cypherCursor.Next() {
		entities, err := cypherCursor.GetRow()
		if err != nil {
			panic(err)
		}
		vertex := entities[0].(*age.Vertex)
		fmt.Println(vertex.Id(), vertex.Label(), vertex.Props())
	}

	// Create Paths (Edges)
	_, err = age.ExecCypher(tx, graphName, 0, "MATCH (a:Person), (b:Person) WHERE a.name='%s' AND b.name='%s' CREATE (a)-[r:workWith {weight: %d}]->(b)", "Jack", "Joe", 3)
	if err != nil {
		panic(err)
	}

	_, err = age.ExecCypher(tx, graphName, 0, "MATCH (a:Person {name: '%s'}), (b:Person {name: '%s'}) CREATE (a)-[r:workWith {weight: %d}]->(b)", "Joe", "Andy", 7)
	if err != nil {
		panic(err)
	}

	tx.Commit()

	tx, err = db.Begin()
	if err != nil {
		panic(err)
	}
	// Query Paths with Cypher
	cypherCursor, err = age.ExecCypher(tx, graphName, 1, "MATCH p=()-[:workWith]-() RETURN p")
	if err != nil {
		panic(err)
	}

	for cypherCursor.Next() {
		entities, err := cypherCursor.GetRow()
		if err != nil {
			panic(err)
		}

		path := entities[0].(*age.Path)
		vertexStart := path.GetAsVertex(0)
		edge := path.GetAsEdge(1)
		vertexEnd := path.GetAsVertex(2)

		fmt.Println(vertexStart, edge, vertexEnd)
	}

	// Delete Vertices
	_, err = age.ExecCypher(tx, graphName, 0, "MATCH (n:Person) DETACH DELETE n RETURN *")
	if err != nil {
		panic(err)
	}
	tx.Commit()
}