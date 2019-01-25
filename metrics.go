package metrics

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kshvakov/clickhouse"
)

var (
	dateRanges = map[string]string{"m": "toIntervalMinute", "h": "toIntervalHour"}
	actions    = map[string]string{"count": "count", "sum": "sum", "any": "any", "anyLast": "anyLast", "min": "min", "max": "max", "avg": "avg", "uniq": "uniq", "uniqhll": "uniqHLL12", "median": "median", "varsamp": "varSamp", "stddevsamp": "stddevSamp", "argmin": "argMin"}

	errActionIsNotFound = errors.New("action is not found")
)

// Metric defines structure for metrics representation
type Metric struct {
	Entity   string    `json:"entity"`
	Names    []string  `json:"names"`
	Values   []float32 `json:"values"`
	DateTime time.Time `json:"datetime"`
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
			datetime DateTime,
			names Array(String),
			values Array(Float32),
			d Date MATERIALIZED toDate(datetime)
		) engine=MergeTree(d, datetime, 8192)
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
	m.DateTime = time.Now()
	tx, err := c.client.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (datetime, names, values, entity) VALUES (?, ?, ?, ?)", c.config.DBName))
	if err != nil {
		return fmt.Errorf("unable to prepare transaction: %v", err)
	}
	_, err = stmt.Exec(time.Now(), m.Names, m.Values, m.Entity)
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
			values   []float32
			names    []string
			entity   string
			datetime time.Time
		)
		if err := rows.Scan(&values, &names, &entity, &datetime); err != nil {
			return nil, fmt.Errorf("unable to scan values: %v", err)
		}
		metrics = append(metrics, &Metric{
			DateTime: datetime,
			Values:   values,
			Names:    names,
			Entity:   entity,
		})
	}
	return metrics, nil
}

// Client returns current ClickHouse client
func (c *ClickHouseMetrics) Client() *sql.DB {
	return c.client
}

// QueryByMetric retruns records by the metric name
func (c *ClickHouseMetrics) QueryByMetric(q *Query) ([]interface{}, error) {

	queryReq := ""
	if q.Entity != "" && q.Label != "" {
		queryReq = fmt.Sprintf("SELECT datetime, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s'", q.Label, q.Label, c.config.DBName, q.Entity)
	}
	if q.TsEqual != 0 {
		queryReq = fmt.Sprintf("SELECT datetime, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s' AND ts = %d", q.Label, q.Label, c.config.DBName, q.Entity, q.TsEqual)
	}
	if q.TsGreater > 0 && q.TsLess > 0 {
		queryReq = fmt.Sprintf("SELECT toUInt64(datetime) AS ts, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s' AND ts > %d AND ts < %d", q.Label, q.Label, c.config.DBName, q.Entity, q.TsGreater, q.TsLess)
	}
	if q.Range != "" {
		queryReq = fmt.Sprintf("SELECT datetime, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s' AND datetime > (%s)", q.Label, q.Label, c.config.DBName, q.Entity, constructDateRange(q.Range))
	}
	rows, err := c.client.Query(queryReq)
	if err != nil {
		return nil, fmt.Errorf("unable to apply query: %v", err)
	}
	defer rows.Close()
	metrics := []interface{}{}
	for rows.Next() {
		var (
			entity string
			ts     interface{}
			value  float32
		)
		if err := rows.Scan(&ts, &entity, &value); err != nil {
			return nil, fmt.Errorf("unable to scan values: %v", err)
		}
		metrics = append(metrics, map[string]interface{}{
			"entity": q.Entity,
			q.Label:  value,
			"ts":     ts,
		})
	}
	return metrics, nil
}

// Aggregate provides operations for aggregation
func (c *ClickHouseMetrics) Aggregate(q *AggregateQuery) (interface{}, error) {
	action, err := checkAction(q.Action)
	if err != nil {
		return nil, err
	}
	queryReq := fmt.Sprintf("SELECT %s(values[indexOf(names, '%s')]) AS result FROM %s", action, q.Label, c.config.DBName)
	if len(q.Entities) > 0 {
		queryReq += q.makeEntitiesQuery()
	}
	fmt.Println(queryReq)
	rows, err := c.client.Query(queryReq)
	if err != nil {
		return nil, fmt.Errorf("unable to apply query: %v", err)
	}
	defer rows.Close()
	var result interface{}
	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			return nil, fmt.Errorf("unable to scan values: %v", err)
		}
	}
	return result, nil
}

// constructDateRange provides constructing of the range
// to ClickHouse format
func constructDateRange(r string) string {
	resp := "now()"
	for k, v := range dateRanges {
		if strings.HasSuffix(r, k) {
			value := r[:len(r)-1]
			return resp + fmt.Sprintf(" - %s(%s)", v, value)
		}
	}
	return resp
}
