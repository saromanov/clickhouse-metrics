package metrics

import (
	"fmt"
	"strings"
)

// queryBuilder provides making of the query to ClickHouse
type queryBuilder struct {
	aq *AggregateQuery
	c  *Config
	q  string
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
	if q.aq.Range != "" {
		queryReq += q.makeRangeQuery()
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
	q.q = res + ")"
	return q.q
}

// makeRangeQuery retruns query if range is defined
func (q *queryBuilder) makeRangeQuery() string {
	if strings.Contains(q.q, "WHERE") {
		q.q += fmt.Sprintf(" AND datetime > (%s)", constructDateRange(q.aq.Range))
		return q.q
	}
	q.q += fmt.Sprintf("WHERE datetime > (%s)", constructDateRange(q.aq.Range))
	return q.q
}

// checkAction return error if action function is not defined
func (q *queryBuilder) checkAction() (string, error) {
	res, ok := actions[q.aq.Action]
	if !ok {
		return "", errActionIsNotFound
	}
	return res, nil
}

// constructDateRange provides constructing of the range
// to ClickHouse format
func constructDateRange(r string) string {
	resp := "now()"
	for k, v := range dateRanges {
		if strings.HasSuffix(r, k) {
			value := r[:len(r)-1]
			return resp + fmt.Sprintf(" - %s(%s)", v, value)
		}
	}
	return resp
}
