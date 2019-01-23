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
func New() (*ClickHouseMetrics, error) {
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?username=&compress=true&debug=true")
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			ts UInt64,
			names Array(String),
			values Array(String),
		) engine=MergeTree(d, timestamp, 8192)
	`)
	if err != nil {
		return nil, fmt.Errorf("unable to create metrics table: %v", err)
	}

	return &ClickHouseMetrics{
		client: connect,
	}, nil
}
