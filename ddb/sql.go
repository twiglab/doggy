package ddb

import "fmt"

const (
	sqlQuery  = "SELECT sn, uuid, code, x, y, z  FROM  "
	createSql = "create or replace table %s as %s "
)

func querySql(tbl string) string {
	return sqlQuery + tbl + " where uuid = $1"
}

func listSql(tbl string) string {
	return sqlQuery + tbl
}

func nextTbl(curr string) string {
	if curr == "cmdb_a" {
		return "cmdb_b"
	}

	return "cmdb_a"
}

func losdSql(t, load string) (tbl string, create string) {
	tbl = nextTbl(t)
	create = fmt.Sprintf(createSql, tbl, load)
	return
}
