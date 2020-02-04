package main

import (
	"math/rand"
	"time"

	metrics "github.com/saromanov/clickhouse-metrics"
)

func main() {
	d, err := metrics.New(&metrics.Config{
		DBName: "base3",
	})
	if err != nil {
		panic(err)
	}

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
}
