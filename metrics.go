package metrics

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kshvakov/clickhouse"
)

// Metric defines structure for metrics representation
type Metric struct {
	Entity    string    `json:"entity"`
	Names     []string  `json:"names"`
	Values    []float32 `json:"values"`
	Timestamp uint64    `json:"timestamp"`
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
			entity String,
			ts UInt64,
			names Array(String),
			values Array(Float32),
			d Date MATERIALIZED toDate(round(ts/1000))
		) engine=MergeTree(d, ts, 8192)
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
	m.Timestamp = uint64(time.Now().Unix())
	tx, err := c.client.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (ts, names, values, entity) VALUES (?, ?, ?, ?)", c.config.DBName))
	if err != nil {
		return fmt.Errorf("unable to prepare transaction: %v", err)
	}
	_, err = stmt.Exec(time.Now().Unix(), m.Names, m.Values, m.Entity)
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
			values []float32
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

// QueryByMetric retruns records by the metric name
func (c *ClickHouseMetrics) QueryByMetric(entity, name string) ([]interface{}, error) {
	rows, err := c.client.Query(fmt.Sprintf("SELECT ts, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s'", name, name, c.config.DBName, entity))
	if err != nil {
		return nil, fmt.Errorf("unable to apply query: %v", err)
	}
	defer rows.Close()
	metrics := []interface{}{}
	for rows.Next() {
		var (
			entity string
			ts     uint64
			value  float32
		)
		if err := rows.Scan(&ts, &entity, &value); err != nil {
			return nil, fmt.Errorf("unable to scan values: %v", err)
		}
		metrics = append(metrics, map[string]interface{}{
			"entity": entity,
			name:     value,
			"ts":     ts,
		})
	}
	return metrics, nil
}
