package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
)

func GetReady(db *sql.DB, graphName string) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("LOAD 'age';")
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("SET search_path = ag_catalog, '$user', public;")
	if err != nil {
		return false, err
	}

	var count int = 0

	err = tx.QueryRow("SELECT count(*) FROM ag_graph WHERE name=$1", graphName).Scan(&count)

	if err != nil {
		return false, err
	}

	if count == 0 {
		_, err = tx.Exec("SELECT create_graph($1);", graphName)
		if err != nil {
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}

type CursorProvider func(columnCount int, rows *sql.Rows) Cursor

type Cursor interface {
	Next() bool
	Close() error
}

func execCypher(cursorProvider CursorProvider, tx *sql.Tx, graphName string, columnCount int, cypher string, args ...interface{}) (Cursor, error) {
	var buf bytes.Buffer

	cypherStmt := fmt.Sprintf(cypher, args...)

	buf.WriteString("SELECT * from cypher('")
	buf.WriteString(graphName)
	buf.WriteString("', $$ ")
	buf.WriteString(cypherStmt)
	buf.WriteString(" $$)")
	buf.WriteString(" as (")
	buf.WriteString("v0 agtype")
	for i := 1; i < columnCount; i++ {
		buf.WriteString(fmt.Sprintf(", v%d agtype", i))
	}
	buf.WriteString(")")

	stmt := buf.String()

	if columnCount == 0 {
		_, err := tx.Exec(stmt)
		if err != nil {
			fmt.Println(stmt)
			return nil, err
		}
		return nil, nil
	} else {
		rows, err := tx.Query(stmt)
		if err != nil {
			fmt.Println(stmt)
			return nil, err
		}
		return cursorProvider(columnCount, rows), nil
	}
}

// ExecCypher
// CREATE , DROP ....
// MATCH .... RETURN ....
// CREATE , DROP .... RETURN ...
func ExecCypher(tx *sql.Tx, graphName string, columnCount int, cypher string, args ...interface{}) (*CypherCursor, error) {
	cursor, err := execCypher(NewCypherCursor, tx, graphName, columnCount, cypher, args...)
	var cypherCursor *CypherCursor
	if cursor != nil {
		cypherCursor = cursor.(*CypherCursor)
	}
	return cypherCursor, err
}

// ExecCypherMap
// CREATE , DROP ....
// MATCH .... RETURN ....
// CREATE , DROP .... RETURN ...
func ExecCypherMap(tx *sql.Tx, graphName string, columnCount int, cypher string, args ...interface{}) (*CypherMapCursor, error) {
	cursor, err := execCypher(NewCypherMapCursor, tx, graphName, columnCount, cypher, args...)
	var cypherMapCursor *CypherMapCursor
	if cursor != nil {
		cypherMapCursor = cursor.(*CypherMapCursor)
	}
	return cypherMapCursor, err
}

// // ExecCypher execute without return
// // CREATE , DROP .... */
// func ExecCypher2(tx *sql.Tx, graphName string, cypher string, args ...interface{}) error {
// 	cypherStmt := fmt.Sprintf(cypher, args...)
// 	stmt := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
// 		graphName, cypherStmt)

// 	_, err := tx.Exec(stmt)
// 	return err
// }

// // QueryCypher execute with return
// // MATCH .... RETURN ....
// // CREATE , DROP .... RETURN ...
// func QueryCypher(tx *sql.Tx, graphName string, cypher string, args ...interface{}) (*CypherCursor, error) {
// 	cypherStmt := fmt.Sprintf(cypher, args...)
// 	stmt := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
// 		graphName, cypherStmt)

// 	rows, err := tx.Query(stmt)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		return NewCypherCursor(1, rows), nil
// 	}
// }

type Age struct {
	db        *sql.DB
	graphName string
}

type AgeTx struct {
	age *Age
	tx  *sql.Tx
}

/**
@param dsn host=127.0.0.1 port=5432 dbname=postgres user=postgres password=agens sslmode=disable
*/
func ConnectAge(graphName string, dsn string) (*Age, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	age := &Age{db: db, graphName: graphName}
	_, err = age.GetReady()

	if err != nil {
		db.Close()
		age = nil
	}

	return age, err
}

func NewAge(graphName string, db *sql.DB) *Age {
	return &Age{db: db, graphName: graphName}
}

func (age *Age) GetReady() (bool, error) {
	tx, err := age.db.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("LOAD 'age';")
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("SET search_path = ag_catalog, '$user', public;")
	if err != nil {
		return false, err
	}

	var count int = 0

	err = tx.QueryRow("SELECT count(*) FROM ag_graph WHERE name=$1", age.graphName).Scan(&count)

	if err != nil {
		return false, err
	}

	if count == 0 {
		_, err = tx.Exec("SELECT create_graph($1);", age.graphName)
		if err != nil {
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}

func (a *Age) Close() error {
	return a.db.Close()
}

func (a *Age) DB() *sql.DB {
	return a.db
}

func (a *Age) Begin() (*AgeTx, error) {
	ageTx := &AgeTx{age: a}
	tx, err := a.db.Begin()
	if err != nil {
		return nil, err
	}
	ageTx.tx = tx
	return ageTx, err
}

func (t *AgeTx) Commit() error {
	return t.tx.Commit()
}

func (t *AgeTx) Rollback() error {
	return t.tx.Rollback()
}

// func buildCypher(graphName string, cypherStmt string, columns )

// def buildCypher(graphName:str, cypherStmt:str, columns:list) ->str:
//     if graphName == None:
//         raise _EXCEPTION_GraphNotSet

//     columnExp=[]
//     if columns != None and len(columns) > 0:
//         for col in columns:
//             if col.strip() == '':
//                 continue
//             elif WHITESPACE.search(col) != None:
//                 columnExp.append(col)
//             else:
//                 columnExp.append(col + " agtype")
//     else:
//         columnExp.append('v agtype')

//     stmtArr = []
//     stmtArr.append("SELECT * from cypher('")
//     stmtArr.append(graphName)
//     stmtArr.append("', $$ ")
//     stmtArr.append(cypherStmt)
//     stmtArr.append(" $$) as (")
//     stmtArr.append(','.join(columnExp))
//     stmtArr.append(");")
//     return "".join(stmtArr)

/** CREATE , DROP .... */
func (a *AgeTx) ExecCypher(columnCount int, cypher string, args ...interface{}) (*CypherCursor, error) {
	return ExecCypher(a.tx, a.age.graphName, columnCount, cypher, args...)
}

// /** CREATE , DROP .... */
// func (a *AgeTx) ExecCypher2(cypher string, args ...interface{}) error {
// 	cypherStmt := fmt.Sprintf(cypher, args...)
// 	stmt := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
// 		a.age.graphName, cypherStmt)

// 	_, err := a.tx.Exec(stmt)
// 	return err
// }

// /** MATCH .... RETURN .... */
// func (a *AgeTx) QueryCypher(cypher string, args ...interface{}) (Cursor, error) {
// 	cypherStmt := fmt.Sprintf(cypher, args...)
// 	stmt := fmt.Sprintf("SELECT * from cypher('%s', $$ %s $$) as (v agtype);",
// 		a.age.graphName, cypherStmt)

// 	rows, err := a.tx.Query(stmt)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		return NewCypherCursor(1, rows), nil
// 	}
// }

type CypherCursor struct {
	Cursor
	columnCount int
	rows        *sql.Rows
	unmarshaler Unmarshaller
}

func NewCypherCursor(columnCount int, rows *sql.Rows) Cursor {
	return &CypherCursor{columnCount: columnCount, rows: rows, unmarshaler: NewAGUnmarshaler()}
}

// func (c *CypherCursor) All() ([]Entity, error) {
// 	defer c.rows.Close()

// 	ens := []Entity{}

// 	for c.rows.Next() {
// 		entity, err := c.GetRow()
// 		if err != nil {
// 			return ens, err
// 		}
// 		ens = append(ens, entity)
// 	}

// 	return ens, nil
// }

func (c *CypherCursor) Next() bool {
	return c.rows.Next()
}

func (c *CypherCursor) GetRow() ([]Entity, error) {
	var gstrs = make([]interface{}, c.columnCount)
	for i := 0; i < c.columnCount; i++ {
		gstrs[i] = new(string)
	}

	err := c.rows.Scan(gstrs...)
	if err != nil {
		return nil, fmt.Errorf("CypherCursor.GetRow:: %s", err)
	}

	entArr := make([]Entity, c.columnCount)
	for i := 0; i < c.columnCount; i++ {
		gstr := gstrs[i].(*string)
		e, err := c.unmarshaler.unmarshal(*gstr)
		if err != nil {
			fmt.Println(i, ">>", gstr)
			return nil, err
		}
		entArr[i] = e
	}

	return entArr, nil
}

func (c *CypherCursor) Close() error {
	return c.rows.Close()
}

//
type CypherMapCursor struct {
	CypherCursor
	mapper *AGMapper
}

func NewCypherMapCursor(columnCount int, rows *sql.Rows) Cursor {
	mapper := NewAGMapper(make(map[string]reflect.Type))
	pcursor := CypherCursor{columnCount: columnCount, rows: rows, unmarshaler: mapper}
	return &CypherMapCursor{CypherCursor: pcursor, mapper: mapper}
}

func (c *CypherMapCursor) PutType(label string, tp reflect.Type) {
	c.mapper.PutType(label, tp)
}

func (c *CypherMapCursor) GetRow() ([]interface{}, error) {
	entities, err := c.CypherCursor.GetRow()

	if err != nil {
		return nil, fmt.Errorf("CypherMapCursor.GetRow:: %s", err)
	}

	elArr := make([]interface{}, c.columnCount)

	for i := 0; i < c.columnCount; i++ {
		ent := entities[i]
		if ent.GType() == G_MAP_PATH {
			elArr[i] = ent
		} else {
			elArr[i] = ent.(*SimpleEntity).Value()
		}
	}

	return elArr, nil
}
