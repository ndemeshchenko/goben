package goben

import (
	"fmt"
	"testing"
	"time"
)

type MockJob001 struct{}

func (d MockJob001) Run() {
	fmt.Println("wassup")
}

func TestDo(*testing.T) {
	var mj1 MockJob001

	g := New()
	g.Every(5).Do(mj1)

	defer g.Cutoff()
	g.Start()

	time.Sleep(5 * time.Second)
}
