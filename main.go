package main

import (
	"github.com/davesavic/go-bench/bench"
	"log"
	"time"
)

func main() {
	bm := bench.NewBenchmark(bench.Options{
		Logger: log.Default(),
	})

	bm.MeasureNParallel(5, func() {
		time.Sleep(1 * time.Second)
	})
}
