package idb

import (
	"slices"
	"strings"
)

const sql = `
select
  sum("in") as "in_total" ,
  sum("out") as "out_total"
from
  count
where
  time >= $start and time < $end
  and device_id in
`

const (
	TIME_START = "start"
	TIME_END   = "end"
	IN_TOTAL   = "in_total"
)

func sqlIn(tables []string) string {
	ss := slices.Clone(tables)
	for i, s := range tables {
		ss[i] = "'" + strings.TrimSpace(s) + "'"
	}
	return "  (" + strings.Join(ss, ",") + ")"
}

func sumOfPointsSQL(tables []string) string {
	inStr := sqlIn(tables)
	return sql + inStr
}
