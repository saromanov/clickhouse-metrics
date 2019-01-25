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

// Query provides struct for the query definition
type Query struct {
	Entities  []string
	Label     string
	TsEqual   uint32
	TsGreater uint32
	TsLess    uint32
	Range     string
}

// GetEntitis returns slice of entities
func (q *Query) GetEntitis() []string {
	return q.Entities
}

// GetLabel returns label
func (q *Query) GetLabel() string {
	return q.Label
}

// GetTsEqual returns timestamp equal param
func (q *Query) GetTsEqual() uint32 {
	return q.TsEqual
}

// GetTsGreater returns timestamp greater param
func (q *Query) GetTsGreater() uint32 {
	return q.TsGreater
}

// GetTsLess returns timestamp Less param
func (q *Query) GetTsLess() uint32 {
	return q.TsLess
}

// GetRange returns range
func (q *Query) GetRange() string {
	return q.Range
}

// GetAction returns action
func (q *Query) GetAction() string {
	return ""
}

// GetActionIf returns action if
func (q *Query) GetActionIf() string {
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
