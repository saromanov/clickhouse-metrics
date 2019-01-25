# clickhouse-metrics
[![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/clickhouse-metrics)](https://goreportcard.com/report/github.com/saromanov/clickhouse-metrics)
[![Coverage Status](https://coveralls.io/repos/github/saromanov/clickhouse-metrics/badge.svg?branch=master)](https://coveralls.io/github/saromanov/clickhouse-metrics?branch=master)
[![CircleCI](https://circleci.com/gh/saromanov/clickhouse-metrics.svg?style=svg)](https://circleci.com/gh/saromanov/clickhouse-metrics)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/f02abbdd25ec4dac9cfb797e1bf2cce7)](https://www.codacy.com/app/saromanov/clickhouse-metrics?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=saromanov/clickhouse-metrics&amp;utm_campaign=Badge_Grade)

Implementation of metric storage over ClickHouse

### Examples - Getting list of the metrics

First step - init of the app

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/saromanov/clickhouse-metrics"
)

func main() {
	d, err := metrics.New(&metrics.Config{
		DBName: "base3",
	})
	if err != nil {
		panic(err)
    }
}
```

Next, add some metric records
```go
s1 := rand.NewSource(time.Now().UnixNano())
r1 := rand.New(s1)
err = d.Insert(&metrics.Metric{
		Entity: "param",
		Names:  []string{"cpu", "load"},
		Values: []float32{float32(r1.Float64()), float32(r1.Float64())},
	})
	if err != nil {
		panic(err)
	}
	err = d.Insert(&metrics.Metric{
		Entity: "foobar",
		Names:  []string{"cpu", "goals"},
		Values: []float32{float32(r1.Float64()), float32(r1.Float64())},
	})
	if err != nil {
		panic(err)
    }
```

And after this, you get apply query by metric which returns list of results which satisfy the conditions

```go
ms, err := d.QueryByMetric(&metrics.Query{
		Label:  "cpu",
		Entity: "param",
		Range:  "1h",
})
```
