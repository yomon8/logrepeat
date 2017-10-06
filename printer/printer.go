package printer

import (
	"fmt"
	"sync"
)

const (
	SpoolSize = 1000
)

type Printer struct {
	Spool chan string
	wg    *sync.WaitGroup
}

var instance *Printer

func Get() *Printer {
	if instance == nil {
		instance = new(Printer)
		instance.Spool = make(chan string, SpoolSize)
		instance.wg = new(sync.WaitGroup)
	}
	return instance
}

func Close() {
	close(instance.Spool)
	instance.wg.Wait()
	fmt.Println("printer closed ..")
}

func (p *Printer) Run() {
	p.wg.Add(1)
	for {
		output, more := <-p.Spool
		if !more {
			break
		}
		fmt.Println(output)
	}
	p.wg.Done()
}
