package metrics

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errQueryIsNotDefined = errors.New("query is not defined")
	errLabelIsNotDefined = errors.New("label is not defined")
)

// queryBuilder provides making of the query to ClickHouse
type queryBuilder struct {
	aq query
	c  *Config
	q  string
}

// make retruns query for ClickHouse
func (q *queryBuilder) make() (string, error) {
	if err := q.validateQuery(); err != nil {
		return "", err
	}
	action, err := q.checkAction()
	if err != nil {
		return "", err
	}
	queryReq := fmt.Sprintf("SELECT %s(values[indexOf(names, '%s')]) AS result FROM %s", action, q.aq.GetLabel(), q.c.DBName)
	if len(q.aq.GetEntities()) > 0 {
		queryReq += q.makeEntitiesQuery()
	}
	if q.aq.GetRange() != "" {
		queryReq += q.makeRangeQuery()
	}
	return queryReq, nil
}

// validateQuery provides validation of the query
func (q *queryBuilder) validateQuery() error {
	if q.aq == nil {
		return errQueryIsNotDefined
	}
	if q.aq.GetLabel() == "" {
		return errLabelIsNotDefined
	}
	return nil
}

func (q *queryBuilder) makeEntitiesQuery() string {
	entities := q.aq.GetEntities()
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
	var res string
	defer func(value string) {
		q.q += res
	}(res)
	if strings.Contains(q.q, "WHERE") {
		res = fmt.Sprintf(" AND (datetime > (%s))", constructDateRange(q.aq.GetRange()))
		return res
	}
	res = fmt.Sprintf("WHERE datetime > (%s)", constructDateRange(q.aq.GetRange()))
	return res
}

// checkAction return error if action function is not defined
func (q *queryBuilder) checkAction() (string, error) {
	res, ok := actions[q.aq.GetAction()]
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
