package preset

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ProgressTracker struct {
	total     int64
	processed int64
	start     time.Time
	active    int32
	mu        sync.Mutex
	lastPrint time.Time
}

func NewProgressTracker(total int) *ProgressTracker {
	return &ProgressTracker{
		total: int64(total),
		start: time.Now(),
	}
}

func (p *ProgressTracker) StartJob() {
	atomic.AddInt32(&p.active, 1)
}

func (p *ProgressTracker) CompleteJob() {
	atomic.AddInt64(&p.processed, 1)
	atomic.AddInt32(&p.active, -1)
	p.maybePrint()
}

func (p *ProgressTracker) maybePrint() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if now.Sub(p.lastPrint) > 100*time.Millisecond {
		p.print()
		p.lastPrint = now
	}
}

func (p *ProgressTracker) print() {
	processed := atomic.LoadInt64(&p.processed)
	active := atomic.LoadInt32(&p.active)
	elapsed := time.Since(p.start)

	percent := float64(processed) / float64(p.total) * 100
	rate := 0.0
	if elapsed.Seconds() > 0 {
		rate = float64(processed) / elapsed.Seconds()
	}

	fmt.Printf("\r[%3.0f%%] %d/%d files | active:%d | %.1f files/sec | %v",
		percent, processed, p.total, active, rate, elapsed.Round(time.Millisecond))
}

func (p *ProgressTracker) Finish() {
	// Завершающая печать
	processed := atomic.LoadInt64(&p.processed)
	elapsed := time.Since(p.start)
	rate := float64(processed) / elapsed.Seconds()

	fmt.Printf("\r[100%%] %d/%d files | %.1f files/sec | %v\n",
		p.total, p.total, rate, elapsed.Round(time.Millisecond))
}
