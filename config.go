package metrics

// Config applies specification to the app start
type Config struct {
	// DBName is a name for metrics database
	DBName string
	// Address is a connection address to ClickHouse
	Address string
}
