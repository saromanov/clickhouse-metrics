# clickhouse-metrics
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
