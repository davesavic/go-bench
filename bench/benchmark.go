package bench

import (
	"sync"
	"time"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Options struct {
	Logger Logger
}

type Benchmark struct {
	Options   Options
	Durations map[int]time.Duration
	Average   time.Duration
}

func NewBenchmark(opts Options) *Benchmark {
	return &Benchmark{
		Durations: map[int]time.Duration{},
		Options:   opts,
	}
}

func (b *Benchmark) Clear() {
	b.Durations = map[int]time.Duration{}
	b.Average = 0
}

func (b *Benchmark) Measure(f func()) {
	b.Clear()
	start := time.Now()
	f()
	b.Durations[0] = time.Since(start)
	b.Average = b.Durations[0]

	if b.Options.Logger != nil {
		b.Options.Logger.Printf("duration: %v", b.Durations[0])
	}
}

func (b *Benchmark) MeasureN(n int, f func()) {
	avg := time.Duration(0)
	for i := 0; i < n; i++ {
		start := time.Now()
		f()
		b.Durations[i] = time.Since(start)
		avg += b.Durations[i]
	}
	b.Average = avg / time.Duration(n)

	if b.Options.Logger != nil {
		for i, duration := range b.Durations {
			b.Options.Logger.Printf("duration %v: %v\n", i, duration)
		}

		b.Options.Logger.Printf("average: %v", b.Average)
	}
}

func (b *Benchmark) MeasureNParallel(n int, f func()) {
	wg := sync.WaitGroup{}
	wg.Add(n)

	rw := sync.RWMutex{}

	avg := time.Duration(0)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			start := time.Now()
			f()

			rw.Lock()
			b.Durations[i] = time.Since(start)
			rw.Unlock()
			avg += b.Durations[i]
			wg.Done()
		}()
	}
	wg.Wait()
	b.Average = avg / time.Duration(n)

	if b.Options.Logger != nil {
		for i, duration := range b.Durations {
			b.Options.Logger.Printf("duration %v: %v\n", i, duration)
		}

		b.Options.Logger.Printf("average: %v", b.Average)
	}
}
