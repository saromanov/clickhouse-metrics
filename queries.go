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
