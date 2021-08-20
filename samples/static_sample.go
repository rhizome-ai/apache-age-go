package samples

import (
	"database/sql"
	"fmt"

	age "github.com/rhizome-ai/apache-age-go"
)

func main() {
	var dsn string = "host={host} port={port} dbname={dbname} user={username} password={password} sslmode=disable"
	var graphName string = "{graph_path}"

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Confirm graph_path created
	_, err = age.GetReady(db, graphName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Tx begin for execute create vertex
	cursor, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create vertices with Cypher
	_, err = age.ExecCypher(cursor, graphName, 0, "CREATE (n:Person {name: '%s', weight:%f})", "Joe", 67.3)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = age.ExecCypher(cursor, graphName, 0, "CREATE (n:Person {name: '%s', weight:77.3, roles:['Dev','marketing']})", "Jack")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = age.ExecCypher(cursor, graphName, 0, "CREATE (n:Person {name: '%s', weight:%d})", "Andy", 59)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Commit Tx
	cursor.Commit()

	// Tx begin for queries
	cursor, err = db.Begin()

	// Query cypher
	cypherCursor, err := age.ExecCypher(cursor, graphName, 1, "MATCH (n:Person) RETURN n")

	// Unmarsal result data to Vertex row by row
	for cypherCursor.Next() {
		entities, err := cypherCursor.GetRow()
		if err != nil {
			fmt.Println(err)
			return
		}
		vertex := entities[0].(*age.Vertex)
		fmt.Println(vertex.Id(), vertex.Label(), vertex.Props())
	}

	// Create Paths (Edges)
	_, err = age.ExecCypher(cursor, graphName, 0, "MATCH (a:Person), (b:Person) WHERE a.name='%s' AND b.name='%s' CREATE (a)-[r:workWith {weight: %d}]->(b)", "Jack", "Joe", 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = age.ExecCypher(cursor, graphName, 0, "MATCH (a:Person {name: '%s'}), (b:Person {name: '%s'}) CREATE (a)-[r:workWith {weight: %d}]->(b)", "Joe", "Andy", 7)
	if err != nil {
		fmt.Println(err)
		return
	}

	cursor.Commit()

	cursor, err = db.Begin()

	// Query Paths with Cypher
	cypherCursor, err = age.ExecCypher(cursor, graphName, 1, "MATCH p=()-[:workWith]-() RETURN p")
	if err != nil {
		fmt.Println(err)
		return
	}

	for cypherCursor.Next() {
		entities, err := cypherCursor.GetRow()
		if err != nil {
			fmt.Println(err)
			return
		}

		path := entities[0].(*age.Path)
		vertexStart := path.GetAsVertex(0)
		edge := path.GetAsEdge(1)
		vertexEnd := path.GetAsVertex(2)

		fmt.Println(vertexStart, edge, vertexEnd)
	}

	// Delete Vertices
	_, err = age.ExecCypher(cursor, graphName, 0, "MATCH (n:Person) DETACH DELETE n RETURN *")
	if err != nil {
		fmt.Println(err)
		return
	}
	cursor.Commit()
}
