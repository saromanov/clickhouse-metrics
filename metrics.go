package metrics

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kshvakov/clickhouse"
)

// Metric defines structure for metrics representation
type Metric struct {
	Entity    string   `json:"entity"`
	Names     []string `json:"names"`
	Values    []string `json:"values"`
	Timestamp uint64   `json:"timestamp"`
}

// ClickHouseMetrics implements the main app
type ClickHouseMetrics struct {
	client *sql.DB
	config *Config
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
		config: c,
	}, nil
}

// Insert provides inserting of the metrics data
func (c *ClickHouseMetrics) Insert(m *Metric) error {
	tx, err := c.client.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (ts, names, values) VALUES (?, ?, ?)", c.config.DBName))
	if err != nil {
		return fmt.Errorf("unable to prepare transaction: %v", err)
	}
	_, err = stmt.Exec(time.Now().Unix(), m.Names, m.Values)
	if err != nil {
		return fmt.Errorf("unable to apply query: %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %v", err)
	}
	return nil
}

// Query provides query of the data by the request
func (c *ClickHouseMetrics) Query(q string) ([]*Metric, error) {
	rows, err := c.client.Query(q)
	if err != nil {
		return nil, fmt.Errorf("unable to apply query: %v", err)
	}
	defer rows.Close()
	metrics := []*Metric{}
	for rows.Next() {
		var (
			values []string
			names  []string
			entity string
			ts     uint64
		)
		if err := rows.Scan(&values, &names, &entity, &ts); err != nil {
			return nil, fmt.Errorf("unable to scan values: %v", err)
		}
		metrics = append(metrics, &Metric{
			Timestamp: ts,
			Values:    values,
			Names:     names,
			Entity:    entity,
		})
	}
	return metrics, nil
}
