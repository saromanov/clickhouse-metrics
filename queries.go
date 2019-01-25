package metrics

import "fmt"

// Query provides struct for the query definition
type Query struct {
	Entity    string
	Label     string
	TsEqual   uint32
	TsGreater uint32
	TsLess    uint32
	Range     string
}

// AggregateQuery defines struct for making aggregation
type AggregateQuery struct {
	Action   string
	Entities []string
	Label    string
	Range    string
	q        string
}

// makeEntitiesQuery retruns query for entities
func (a *AggregateQuery) makeEntitiesQuery() string {
	if len(a.Entities) == 1 {
		return fmt.Sprintf(" WHERE entity = '%s'", a.Entities[0])
	}
	res := " WHERE ( "
	for i := 0; i < len(a.Entities); i++ {
		res += fmt.Sprintf("entity = '%s' ", a.Entities[i])
		if i+1 != len(a.Entities) {
			res += "OR "
		}
	}
	return res + ")"
}
