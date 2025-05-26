package criteria

import (
	"fmt"
	"strings"
)

func Build(query string, criteria *Criteria, startFrom uint64, valueOperator string) string {
	sql := criteria.Condition.ToSQL()

	fields := strings.Fields(sql)

	for i, f := range fields {
		if f == "?" {
			fields[i] = fmt.Sprintf(valueOperator, startFrom)
			startFrom++
		}
	}

	return fmt.Sprintf("%s\n%s", query, strings.Join(fields, " "))
}