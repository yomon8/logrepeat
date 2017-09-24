package reperter

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/yomon8/logrepeat/printer"
	"github.com/yomon8/logrepeat/request"
)

const (
	requestsBufferSize = 1000
	resultBufferSize   = 1000
	logIntervalSec     = 3
)

var (
	resultBuffer chan *Result = make(chan *Result, resultBufferSize)
)

// Repeater run repeat request generated from original requests
type Repeater struct {
	requests *request.Requests
	wg       *sync.WaitGroup
	buffer   chan *request.Request
	quit     chan bool
	count    int
	total    int
}

// NewRepeater create Repeater instanse
func NewRepeater(requests *request.Requests) *Repeater {
	r := new(Repeater)
	r.requests = requests
	r.wg = new(sync.WaitGroup)
	r.buffer = make(chan *request.Request, requestsBufferSize)
	r.quit = make(chan bool)
	r.total = requests.Len()
	return r
}

func (r *Repeater) collectStats() {
	var statsTime, laststatsTime int64
	var results Results = make([]*Result, 0)
	for {
		select {
		case <-r.quit:
			printer.Get().Spool <- fmt.Sprint("analyzer stopped ...")
			return
		case result := <-resultBuffer:
			results = append(results, result)
		default:
			statsTime = time.Now().Unix()
			if statsTime-laststatsTime > logIntervalSec && len(results) > 0 {
				r.count += len(results)
				progress := float32(r.count) / float32(r.total) * 100
				sort.Sort(results)
				printer.Get().Spool <- fmt.Sprintf("%s - %s  %s  (%.1f%%Done)",
					results[0].requestTimeString(),
					results[len(results)-1].requestTimeString(),
					results.GetStatsString(),
					progress)

				results = make([]*Result, 0)
				laststatsTime = statsTime
			}
		}
	}
	time.Sleep(10 * time.Microsecond)
}

func (r *Repeater) request(isDryrun bool) {
	r.wg.Add(1)
	for {
		req, more := <-r.buffer
		if !more {
			printer.Get().Spool <- fmt.Sprint("worker stopped ...")
			break
		}
		var code int
		var responseTime time.Duration
		for {
			repeatTime := time.Now()
			if repeatTime.Sub(req.RepeatTime) <= 0 {
				time.Sleep(10 * time.Microsecond)
				continue
			} else {
				if isDryrun {
					code = 999
					responseTime = time.Duration(0)
				} else {
					client := &http.Client{Timeout: time.Duration(20) * time.Second}
					httpreq, err := http.NewRequest("GET", req.Url, nil)
					if err != nil {
						code = 1000
						responseTime = time.Duration(0)
					}
					start := time.Now()
					if res, err := client.Do(httpreq); err != nil {
						code = 1001
					} else {
						code = res.StatusCode
					}
					responseTime = time.Now().Sub(start)
				}
				resultBuffer <- newResult(repeatTime, code, responseTime)
				break
			}
		}
	}
	r.wg.Done()
}

// Run all repeat requests
func (r *Repeater) Run(concurrency int, isDryrun bool) {
	go printer.Get().Run()

	for i := 0; concurrency > i; i++ {
		go r.request(isDryrun)
	}

	go r.collectStats()

	for _, req := range *r.requests {
		r.buffer <- req
	}
	close(r.buffer)
	r.wg.Wait()
	close(resultBuffer)
	r.quit <- true
	printer.Close()
	color.Blue("Repeat Completed!")
}
