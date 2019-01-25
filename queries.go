package metrics

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
}
