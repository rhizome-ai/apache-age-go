# Apache AGE Python Driver

[Apache AGE](https://age.apache.org/) is a PostgreSQL extension that provides graph database functionality. The goal of the Apache AGE project is to create single storage that can handle both relational and graph model data so that users can use standard ANSI SQL along with openCypher, the Graph query language. This repository hosts the development of the Python driver for this Apache extension (currently in Incubator status). Thanks for checking it out.

A graph consists of a set of vertices (also called nodes) and edges, where each individual vertex and edge possesses a map of properties. A vertex is the basic object of a graph, that can exist independently of everything else in the graph. An edge creates a directed connection between two vertices. A graph database is simply composed of vertices and edges. This type of database is useful when the meaning is in the relationships between the data. Relational databases can easily handle direct relationships, but indirect relationships are more difficult to deal with in relational databases. A graph database stores relationship information as a first-class entity. Apache AGE gives you the best of both worlds, simultaneously.

Apache AGE is:

- **Powerful** -- AGE adds graph database support to the already popular PostgreSQL database: PostgreSQL is used by organizations including Apple, Spotify, and NASA.
- **Flexible** -- AGE allows you to perform openCypher queries, which make complex queries much easier to write.
- **Intelligent** -- AGE allows you to perform graph queries that are the basis for many next level web services such as fraud & intrustion detection, master data management, product recommendations, identity and relationship management, experience personalization, knowledge management and more.

### Features
* Cypher query support for database/sql PostreSQL driver (enables cypher queries directly)
* Deserialize AGE result (AGType) to Vertex, Edge, Path

### Go get & gomod
* over Go 1.16
* This module runs on golang standard api [database/sql](https://golang.org/pkg/database/sql/) and [antlr4-go](https://github.com/antlr/antlr4/tree/master/runtime/Go/antlr)

```(shell)
go get github.com/rhizome-ai/apache-age-go
```
#### gomod
```(go)
require  github.com/rhizome-ai/apache-age-go {version}
 
```
Check [latest version](https://github.com/rhizome-ai/apache-age-go/releases)

### For more information about [Apache AGE](https://age.apache.org/)
* Apache Incubator Age : https://age.apache.org/
* Github : https://github.com/apache/incubator-age
* Document : https://age.incubator.apache.org/docs/

### Check AGE loaded on your PostgreSQL
Connect to your containerized Postgres instance and then run the following commands:
```(sql)
# psql 
CREATE EXTENSION age;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;
```

### Test
Check out and rewrite DSN in incubator-age/drivers/golang/age_test.go
```
cd apache-age-go
go test -v .

```

### Usage 1: using static function 'ExecCypher'
```
import (
    "fmt"
    "database/sql"
	age "github.com/rhizome-ai/apache-age-go"
)

func main() {
    var dsn string = "host={host} port={port} dbname={dbname} user={username} password={password} sslmode=disable"
    var graphName string = "{graph_path}"

    // Connect to PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

    // Confirm graph_path created
	_, err = age.GetReady(db, graphName)
	if err != nil {
		t.Fatal(err)
	}

    // Tx begin for execute create vertex
	cursor, err := db.Begin()
	if err != nil {
		fmt.Println(err)
        return
	}

    // Create vertices with Cypher
	_, err = ExecCypher(cursor, graphName, 0, "CREATE (n:Person {name: '%s', weight:%f})", "Joe", 67.3)
	if err != nil {
		fmt.Println(err)
        return
	}

	_, err = ExecCypher(cursor, graphName, 0, "CREATE (n:Person {name: '%s', weight:77.3, roles:['Dev','marketing']})", "Jack")
	if err != nil {
		fmt.Println(err)
        return
	}

	_, err = ExecCypher(cursor, graphName, 0, "CREATE (n:Person {name: '%s', weight:%d})", "Andy", 59)
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
			t.Fatal(err)
		}
		vertex := entities[0].(*Vertex)
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
			t.Fatal(err)
		}

		path := entities[0].(*Path)
		vertexStart := path.GetAsVertex(0)
		edge := path.GetAsEdge(1)
		vertexEnd := path.GetAsVertex(2)

		fmt.Println(vertexStart, edge, vertexEnd)
	}

    // Delete Vertices 
	_, err = age.ExecCypher(cursor, graphName, "MATCH (n:Person) DETACH DELETE n RETURN *")
	if err != nil {
		fmt.Println(err)
        return
	}
	cursor.Commit()
}
```
### License
Apache-2.0 License