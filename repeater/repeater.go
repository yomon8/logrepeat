package reperter

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/yomon8/logrepeat/printer"
	"github.com/yomon8/logrepeat/request"
)

const (
	REQUEST_BUFFER_SIZE = 1000
	RESULT_BUFFER_SIZE  = 1000
	LOG_INTERVAL_SEC    = 3
)

var (
	resultBuffer chan *Result = make(chan *Result, RESULT_BUFFER_SIZE)
)

type Repeater struct {
	requests *request.Requests
	wg       *sync.WaitGroup
	buffer   chan *request.Request
	quit     chan bool
	count    int
	total    int
}

func NewRepearter(requests *request.Requests) *Repeater {
	r := new(Repeater)
	r.requests = requests
	r.wg = new(sync.WaitGroup)
	r.buffer = make(chan *request.Request, REQUEST_BUFFER_SIZE)
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
			if statsTime-laststatsTime > LOG_INTERVAL_SEC && len(results) > 0 {
				r.count += len(results)
				progress := float32(r.count) / float32(r.total) * 100
				sort.Sort(results)
				printer.Get().Spool <- fmt.Sprintf("%s - %s  %s  (%.1f%%)",
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
		for {
			repeatTime := time.Now()
			var code int
			if isDryrun {
				code = 999
			} else {
				client := &http.Client{Timeout: time.Duration(20) * time.Second}
				req, err := http.NewRequest("GET", req.Url, nil)
				if err != nil {
					log.Println("request error:", err)
				}
				if res, err := client.Do(req); err != nil {
					log.Println("reqest do error:", err)
				} else {
					code = res.StatusCode
				}
			}
			if repeatTime.Sub(req.RepeatTime) >= 0 {
				resultBuffer <- newResult(repeatTime, code)
				break
			} else {
				time.Sleep(10 * time.Microsecond)
			}
		}
	}
	r.wg.Done()
}

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
