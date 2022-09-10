package common

import (
	"fmt"
	"time"
)

type Timer struct {
	name  string
	start time.Time
}

func Start(name string) Timer {
	fmt.Printf("\n--> Started %s...\n", name)
	return Timer{name, time.Now()}
}

func (t Timer) Complete() {
	fmt.Printf("... %s took %v\n", t.name, time.Since(t.start))
}
