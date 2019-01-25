package metrics

type query interface {
	GetEntities() []string
	GetLabel() string
	GetTsEqual() uint32
	GetTsGreater() uint32
	GetTsLess() uint32
	GetRange() string
	GetAction() string
	GetActionIf() string
}

// ListQuery provides struct for the query definition
type ListQuery struct {
	Entities  []string
	Label     string
	TsEqual   uint32
	TsGreater uint32
	TsLess    uint32
	Range     string
}

// GetEntitis returns slice of entities
func (q *ListQuery) GetEntitis() []string {
	return q.Entities
}

// GetLabel returns label
func (q *ListQuery) GetLabel() string {
	return q.Label
}

// GetTsEqual returns timestamp equal param
func (q *ListQuery) GetTsEqual() uint32 {
	return q.TsEqual
}

// GetTsGreater returns timestamp greater param
func (q *ListQuery) GetTsGreater() uint32 {
	return q.TsGreater
}

// GetTsLess returns timestamp Less param
func (q *ListQuery) GetTsLess() uint32 {
	return q.TsLess
}

// GetRange returns range
func (q *ListQuery) GetRange() string {
	return q.Range
}

// GetAction returns action
func (q *ListQuery) GetAction() string {
	return ""
}

// GetActionIf returns action if
func (q *ListQuery) GetActionIf() string {
	return ""
}

// AggregateQuery defines struct for making aggregation
type AggregateQuery struct {
	Action   string
	Entities []string
	Label    string
	Range    string
	ActionIf string
}

// GetEntitis returns slice of entities
func (q *AggregateQuery) GetEntitis() []string {
	return q.Entities
}

// GetLabel returns label
func (q *AggregateQuery) GetLabel() string {
	return q.Label
}

// GetTsEqual returns timestamp equal param
func (q *AggregateQuery) GetTsEqual() uint32 {
	return 0
}

// GetTsGreater returns timestamp greater param
func (q *AggregateQuery) GetTsGreater() uint32 {
	return 0
}

// GetTsLess returns timestamp Less param
func (q *AggregateQuery) GetTsLess() uint32 {
	return 0
}

// GetRange returns range
func (q *AggregateQuery) GetRange() string {
	return q.Range
}

// GetAction returns action
func (q *AggregateQuery) GetAction() string {
	return q.Action
}

// GetActionIf returns action if
func (q *AggregateQuery) GetActionIf() string {
	return q.ActionIf
}
