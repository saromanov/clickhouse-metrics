package metrics

import (
	"database/sql"
	"fmt"

	"github.com/kshvakov/clickhouse"
)

// ClickHouseMetrics implements the main app
type ClickHouseMetrics struct {
	client *sql.DB
}

// New provides initialization of the project
func New(c *Config) (*ClickHouseMetrics, error) {
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?username=&compress=true&debug=true")
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return nil, fmt.Errorf("[%d] %s %s", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, fmt.Errorf("unable to ping Clickhouse: %v", err)
	}

	if c.DBName == "" {
		c.DBName = "metrics"
	}
	_, err = connect.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			ts UInt64,
			names Array(String),
			values Array(String),
		) engine=MergeTree(d, timestamp, 8192)
	`, c.DBName))
	if err != nil {
		return nil, fmt.Errorf("unable to create metrics table: %v", err)
	}

	return &ClickHouseMetrics{
		client: connect,
	}, nil
}
