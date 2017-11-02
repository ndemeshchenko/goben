package goben

import (
	"fmt"
	"sync"
	"time"
)

// Job interdace definition
type Job interface {
	Run()
}

// Entry of the schedule
type Entry struct {
	Period time.Duration
	Next   time.Time
	Prev   time.Time
	Job    Job
}

func newEntry(period time.Duration) *Entry {
	return &Entry{
		Period: period,
		Prev:   time.Unix(0, 0),
	}
}

// Do is definer of the action
func (e *Entry) Do(job Job) {
	e.Job = job
}

// Goben scheduler
type Goben struct {
	entries []*Entry
	cutoff  chan struct{}
	running bool
}

// New is a factory for Goben
func New() *Goben {
	return &Goben{
		cutoff:  make(chan struct{}),
		entries: nil,
		running: false,
	}
}

func (g *Goben) schedule(e *Entry) {
	e.Prev = time.Now()
	e.Next = e.Prev.Add(e.Period)
}

// Every is a function defining schedule
// returns new Entry pointer
func (g *Goben) Every(duration time.Duration) *Entry {
	// entry := newEntry(time.Duration(seconds) * time.Second)
	fmt.Println("[INFO] Define duration")
	entry := newEntry(duration)

	g.schedule(entry)

	if !g.running {
		g.entries = append(g.entries, entry)
	}

	return entry
}

// Start main entrypoint of the scheduler start
func (g *Goben) Start() {
	if g.running {
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)

	g.running = true
	go g.run(&wg)
	wg.Wait()
}

// Cutoff breaking the execution
func (g *Goben) Cutoff() {
	if !g.running {
		return
	}

	g.cutoff <- struct{}{}
	g.running = false
}

func (g *Goben) run(wg *sync.WaitGroup) {
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				g.RunPending()
				continue
			case <-g.cutoff:
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()

}

// RunPending ...
func (g *Goben) RunPending() {
	go func() {
		for _, entry := range g.entries {
			if time.Now().After(entry.Next) {
				go g.runJob(entry)
			}
		}
	}()
}

func (g *Goben) runJob(e *Entry) {
	defer func() {
		if r := recover(); r != nil {
			g.schedule(e)
		}
	}()

	g.schedule(e)
	e.Job.Run()
}
