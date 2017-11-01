package goben

import "time"

type Job interface {
	Run()
}

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

func (e *Entry) Do(job Job) {
	e.Job = job
}

type Goben struct {
	entries []*Entry
	cutoff  chan struct{}
	running bool
}

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

func (g *Goben) Every(seconds uint64) *Entry {
	entry := newEntry(time.Duration(seconds) * time.Second)

	g.schedule(entry)

	if !g.running {
		g.entries = append(g.entries, entry)
	}

	return entry
}

func (g *Goben) Start() {
	if g.running {
		return
	}
	g.running = true
	go g.run()
}

func (g *Goben) Cutoff() {
	if !g.running {
		return
	}

	g.cutoff <- struct{}{}
	g.running = false
}

func (g *Goben) run() {
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				g.RunPending()
				continue
			case <-g.cutoff:
				ticker.Stop()
				return
			}
		}
	}()

}

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
