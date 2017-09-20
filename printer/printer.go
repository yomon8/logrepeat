package printer

import (
	"fmt"
	"sync"
)

const (
	SPOOL_SIZE = 1000
)

type printer struct {
	Spool chan string
	wg    *sync.WaitGroup
}

var instance *printer

func Get() *printer {
	if instance == nil {
		instance = new(printer)
		instance.Spool = make(chan string, SPOOL_SIZE)
		instance.wg = new(sync.WaitGroup)
	}
	return instance
}

func Close() {
	instance.wg.Wait()
	close(instance.Spool)
}

func (p *printer) Run() {
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
