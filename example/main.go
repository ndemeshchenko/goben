package main

import (
	"fmt"
	"time"

	"github.com/ndemeshchenko/goben"
)

// SampleJob impl for Job interface
type SampleJob struct{}

// Run implementation for SameplJob
func (s SampleJob) Run() {
	fmt.Println(time.Now(), "Each 1/2 seconds")
}

func main() {
	var j SampleJob
	g := goben.New()
	g.Every(500 * time.Millisecond).Do(j)

	defer g.Cutoff()

	go g.Start()

	time.Sleep(5 * time.Second)
}
