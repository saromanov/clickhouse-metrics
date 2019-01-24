package metrics

// Query provides struct for the query definition
type Query struct {
	Entity    string
	Label     string
	TsEqual   uint32
	Tsgreater uint32
	TsLess    uint
	Range     string
}
