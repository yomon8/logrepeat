package printer

import (
	"fmt"
	"sync"
	"time"
)

const (
	SpoolSize = 1000
)

type Printer struct {
	Spool chan string
	wg    *sync.WaitGroup
	warm  bool
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
	go p.worker()
	for !p.warm {
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (p *Printer) worker() {
	p.wg.Add(1)
	for {
		p.warm = true
		output, more := <-p.Spool
		if !more {
			break
		}
		fmt.Println(output)
	}
	p.wg.Done()
}
