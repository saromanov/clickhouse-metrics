package metrics

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"github.com/pkg/errors"

	"github.com/ClickHouse/clickhouse-go"
)

var (
	dateRanges = map[string]string{"m": "toIntervalMinute", "h": "toIntervalHour", "d": "toIntervalDay", "w": "toIntervalWeek"}
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
	connect, err := sql.Open("clickhouse", c.Address)
	if err != nil {
		return nil, err
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return nil, errors.Wrap("[%d] %s %s", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, errors.Wrap("unable to ping Clickhouse", err)
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
		return nil, errors.Wrap("unable to create metrics table", err)
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
		return errors.Wrap(err, "unable to begin transaction")
	}
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (datetime, names, values, entity) VALUES (?, ?, ?, ?)", c.config.DBName))
	if err != nil {
		return errors.Wrap(err, "unable to prepare transaction")
	}
	_, err = stmt.Exec(time.Now(), m.Names, m.Values, m.Entity)
	if err != nil {
		return errors.Wrap(err, "unable to apply query")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit transaction")
	}
	return nil
}

// Query provides query of the data by the request
func (c *ClickHouseMetrics) Query(q string) ([]*Metric, error) {
	rows, err := c.client.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "unable to apply query")
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
			return nil, errors.Wrap(err, "unable to scan values")
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

// List retruns list of the metrics by the query
func (c *ClickHouseMetrics) List(q *ListQuery) ([]interface{}, error) {

	qb := &queryBuilder{
		aq: q,
		c:  c.config,
	}
	query, err := qb.make()
	if err != nil {
		return nil, err
	}

	/*queryReq := ""
	if q.TsEqual != 0 {
		queryReq = fmt.Sprintf("SELECT datetime, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s' AND ts = %d", q.Label, q.Label, c.config.DBName, q.Entity, q.TsEqual)
	}
	if q.TsGreater > 0 && q.TsLess > 0 {
		queryReq = fmt.Sprintf("SELECT toUInt64(datetime) AS ts, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s' AND ts > %d AND ts < %d", q.Label, q.Label, c.config.DBName, q.Entity, q.TsGreater, q.TsLess)
	}
	if q.Range != "" {
		queryReq = fmt.Sprintf("SELECT datetime, entity, values[indexOf(names, '%s')] AS %s FROM %s WHERE entity = '%s' AND datetime > (%s)", q.Label, q.Label, c.config.DBName, q.Entity, constructDateRange(q.Range))
	}*/
	rows, err := c.client.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to apply query")
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
			return nil, errors.Wrap(err, "unable to scan values")
		}
		metrics = append(metrics, map[string]interface{}{
			"entity": entity,
			q.Label:  value,
			"ts":     ts,
		})
	}
	return metrics, nil
}

// Aggregate provides operations for aggregation
func (c *ClickHouseMetrics) Aggregate(q *AggregateQuery) (interface{}, error) {
	qb := &queryBuilder{
		aq: q,
		c:  c.config,
	}
	query, err := qb.make()
	if err != nil {
		return nil, err
	}
	fmt.Println("QUERY: ", query)
	rows, err := c.client.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to apply query")
	}
	defer rows.Close()
	var result interface{}
	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			return nil, errors.Wrap(err, "unable to scan values")
		}
	}
	return result, nil
}
