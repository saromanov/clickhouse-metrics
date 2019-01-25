package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListQuery(t *testing.T) {
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

	qb := &queryBuilder {
		aq: lq,
		c: &Config{
			DBName: "test",
		},
	}
	_, err := qb.make()
	assert.NoError(t, err)
}