package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	var item uint32
	item = 1234567
	lq := &ListQuery{
		Entities:  []string{"one"},
		Label:     "label",
		TsEqual:   item,
		TsGreater: item,
		TsLess:    item,
		Range:     "5m",
	}

	qb := &queryBuilder{
		aq: lq,
		c: &Config{
			DBName: "test",
		},
	}
	q, err := qb.make()
	assert.NoError(t, err)
	assert.Equal(t, q, "SELECT values[indexOf(names, 'label')] AS result FROM test  WHERE entity = 'one'WHERE datetime > (now() - toIntervalMinute(5))", "not equal")

	lq = &ListQuery{
		Entities:  []string{"one", "two"},
		Label:     "label",
		TsEqual:   item,
		TsGreater: item,
		TsLess:    item,
		Range:     "2h",
	}

	qb = &queryBuilder{
		aq: lq,
		c: &Config{
			DBName: "test",
		},
	}
	q, err = qb.make()
	assert.NoError(t, err)
	assert.Equal(t, q, "SELECT values[indexOf(names, 'label')] AS result FROM test  WHERE ( entity = 'one' OR entity = 'two' ) AND (datetime > (now() - toIntervalHour(2)))", "not equal")

	lq = &ListQuery{
		Range: "2h",
	}

	qb = &queryBuilder{
		aq: lq,
		c: &Config{
			DBName: "test",
		},
	}
	_, err = qb.make()
	assert.Error(t, err)
}

func TestInvalidRange(t *testing.T) {
	lq := &ListQuery{
		Entities: []string{"one", "two"},
		Label:    "label",
		Range:    "sdadasdas",
	}

	qb := &queryBuilder{
		aq: lq,
		c: &Config{
			DBName: "test",
		},
	}
	q, err := qb.make()
	assert.NoError(t, err)
	assert.Equal(t, q, "SELECT values[indexOf(names, 'label')] AS result FROM test  WHERE ( entity = 'one' OR entity = 'two' ) AND (datetime > (now() - toIntervalHour(1)))", "not equal")
}

func TestAction(t *testing.T) {
	lq := &AggregateQuery{
		Entities: []string{"one", "two"},
		Label:    "label",
		Range:    "2h",
		Action:   "sum",
	}

	qb := &queryBuilder{
		aq: lq,
		c: &Config{
			DBName: "test",
		},
	}
	q, err := qb.make()
	assert.NoError(t, err)
	assert.Equal(t, q, "SELECT sum(values[indexOf(names, 'label')]) AS result FROM test  WHERE ( entity = 'one' OR entity = 'two' ) AND (datetime > (now() - toIntervalHour(2)))", "not equal")
}

func TestInvalidAction(t *testing.T) {
	lq := &AggregateQuery{
		Entities: []string{"one", "two"},
		Label:    "label",
		Range:    "2h",
		Action:   "sumss",
	}

	qb := &queryBuilder{
		aq: lq,
		c: &Config{
			DBName: "test",
		},
	}
	_, err := qb.make()
	assert.Error(t, err)
}
