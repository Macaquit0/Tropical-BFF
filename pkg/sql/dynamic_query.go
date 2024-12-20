package sql

import (
	"fmt"
)

func AddWhereParams(namedArgs map[string]interface{}, startIndex int) (string, []interface{}) {
	var args []interface{}
	var query string
	if len(namedArgs) > 0 {
		query = "WHERE"
	}
	i := startIndex
	for key, value := range namedArgs {
		if i == startIndex {
			query += fmt.Sprintf(" %s = $%d", key, i)
		} else {
			query += fmt.Sprintf(" AND %s = $%d", key, i)
		}
		args = append(args, value)
		i++
	}

	return query, args
}
