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
	assert.Equal(t, lq.GetEntities(), []string{"one"}, "not equal")
	assert.Equal(t, lq.GetLabel(), "label", "not equal")
	assert.Equal(t, lq.GetTsEqual(), item, "not equal")
	assert.Equal(t, lq.GetTsGreater(), item, "not equal")
	assert.Equal(t, lq.GetTsLess(), item, "not equal")
	assert.Equal(t, lq.GetRange(), "5m", "not equal")
	assert.Equal(t, lq.GetAction(), "", "not equal")
	assert.Equal(t, lq.GetActionIf(), "", "not equal")
}

func TestAggregateQuery(t *testing.T) {
	lq := &AggregateQuery{
		Entities: []string{"one"},
		Label:    "label",
		Action:   "sum",
		ActionIf: "x < 1",
		Range:    "5m",
	}
	assert.Equal(t, lq.GetEntities(), []string{"one"}, "not equal")
	assert.Equal(t, lq.GetLabel(), "label", "not equal")
	assert.Equal(t, lq.GetTsEqual(), uint32(0), "not equal")
	assert.Equal(t, lq.GetTsGreater(), uint32(0), "not equal")
	assert.Equal(t, lq.GetTsLess(), uint32(0), "not equal")
	assert.Equal(t, lq.GetRange(), "5m", "not equal")
	assert.Equal(t, lq.GetAction(), "sum", "not equal")
	assert.Equal(t, lq.GetActionIf(), "x < 1", "not equal")
}
