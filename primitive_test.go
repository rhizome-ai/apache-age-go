package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestDB(t *testing.T) {
	connStr := "host=127.0.0.1 port=5432 dbname=postgres user=postgres password=agens sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("LOAD 'age';")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("SET search_path = ag_catalog, '$user', public;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("SET search_path = ag_catalog, '$user', public;")
	if err != nil {
		log.Fatal(err)
	}

	var count int = 0

	err = db.QueryRow("SELECT count(*) FROM ag_graph WHERE name=$1", "testGraph").Scan(&count)

	if count == 0 {
		_, err = db.Exec("SELECT create_graph($1);", "testGraph")
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
		"testGraph", "CREATE (n:Person {name: 'Joe'})")

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	stmt2 := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
		"testGraph", "MATCH (n:Person) RETURN n")

	rows, err := db.Query(stmt2)

	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var gstr string
			rows.Scan(&gstr)
			fmt.Println(gstr)
		}
	}

	// stmt3 := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
	// 	"testGraph", "MATCH (n:Person {name: $1}) RETURN n")

	// fmt.Println(stmt3)

	// rows, err = db.Query(stmt3, "Joe")

	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	defer rows.Close()
	// 	for rows.Next() {
	// 		var gstr string
	// 		rows.Scan(&gstr)
	// 		fmt.Println(gstr)
	// 	}
	// }

	stmtClear := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
		"testGraph", "MATCH (n:Person) DETACH DELETE n RETURN *")

	_, err = db.Exec(stmtClear)

	if err != nil {
		log.Fatal(err)
	}

}
