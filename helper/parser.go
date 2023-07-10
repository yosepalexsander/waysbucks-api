package helper

import (
	"fmt"
	"regexp"
)

var nameRegex = regexp.MustCompile(`\A[\[\]]*([^\[\]]+)\]*`)
var clauseRegex = regexp.MustCompile(`\[([^\[\]]+)\]`)

func QueryParamsToSqlClauses(queries map[string][]string) ([]string, string) {
	if len(queries) < 1 {
		return nil, ""
	}

	whereClauses := []string{}
	orderByClause := ""

	for k, v := range queries {
		value := v[0]

		if k == "order_by" {
			orderByClause += value
			continue
		}

		field := string(nameRegex.Find([]byte(k)))
		clauses := clauseRegex.FindStringSubmatch(k)

		if len(clauses) < 2 {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", field, value))
			continue
		}

		switch clauses[1] {
		case "gte":
			whereClauses = append(whereClauses, fmt.Sprintf("%s >= %s", field, value))
		case "lte":
			whereClauses = append(whereClauses, fmt.Sprintf("%s <= %s", field, value))
		}
	}

	return whereClauses, orderByClause
}
