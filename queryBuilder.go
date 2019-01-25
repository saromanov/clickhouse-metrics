package metrics

import (
	"fmt"
)

// queryBuilder provides making of the query to ClickHouse
type queryBuilder struct {
	aq *AggregateQuery
	c  *Config
}

// make retruns query for ClickHouse
func (q *queryBuilder) make() (string, error) {
	action, err := q.checkAction()
	if err != nil {
		return "", err
	}
	queryReq := fmt.Sprintf("SELECT %s(values[indexOf(names, '%s')]) AS result FROM %s", action, q.aq.Label, q.c.DBName)
	if len(q.aq.Entities) > 0 {
		queryReq += q.makeEntitiesQuery()
	}
	return queryReq, nil
}

func (q *queryBuilder) makeEntitiesQuery() string {
	entities := q.aq.Entities
	if len(entities) == 1 {
		return fmt.Sprintf(" WHERE entity = '%s'", entities[0])
	}
	res := " WHERE ( "
	for i := 0; i < len(entities); i++ {
		res += fmt.Sprintf("entity = '%s' ", entities[i])
		if i+1 != len(entities) {
			res += "OR "
		}
	}
	return res + ")"
}

// checkAction return error if action function is not defined
func (q *queryBuilder) checkAction() (string, error) {
	res, ok := actions[q.aq.Action]
	if !ok {
		return "", errActionIsNotFound
	}
	return res, nil
}
