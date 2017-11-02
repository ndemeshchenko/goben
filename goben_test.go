package goben

import (
	"fmt"
	"testing"
	"time"
)

type MockJob001 struct{}

func (d MockJob001) Run() {
	fmt.Println("[TEST] ", time.Now(), " Mock job run")
}

func TestDo(*testing.T) {
	var mj1 MockJob001

	g := New()
	g.Every(500 * time.Millisecond).Do(mj1)

	defer g.Cutoff()

	go g.Start()

	time.Sleep(3 * time.Second)
}
