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
	resultBuffer = make(chan *Result, resultBufferSize)
)

// Repeater run repeat request generated from original requests
type Repeater struct {
	requests       *request.Requests
	wg             *sync.WaitGroup
	wgCollectStats *sync.WaitGroup
	buffer         chan *request.RequestEntry
	quit           bool
	count          int
	total          int
}

// NewRepeater create Repeater instanse
func NewRepeater(requests *request.Requests) *Repeater {
	return &Repeater{
		requests:       requests,
		wg:             new(sync.WaitGroup),
		buffer:         make(chan *request.RequestEntry, requestsBufferSize),
		total:          requests.Len(),
		wgCollectStats: new(sync.WaitGroup),
	}
}

func (r *Repeater) collectStats() {
	var statsTime, laststatsTime int64
	var results Results = make([]*Result, 0)
	var printResult func() = func() {
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

	r.wgCollectStats.Add(1)
	for {
		result, more := <-resultBuffer
		if !more {
			if len(results) > 0 {
				printResult()
			}
			printer.Get().Spool <- fmt.Sprint("analyzer stopped ...")
			printer.Close()
			r.wgCollectStats.Done()
			return
		}
		results = append(results, result)
		statsTime = time.Now().Unix()
		if statsTime-laststatsTime > logIntervalSec && len(results) > 0 {
			printResult()
		}
		time.Sleep(10 * time.Microsecond)
	}
}

func (r *Repeater) runRequestWorker(isDryrun bool, ignoreReqTime bool) {
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

			if !ignoreReqTime {
				if repeatTime.Sub(req.RepeatTime) <= 0 {
					//wait for repeat time
					time.Sleep(10 * time.Microsecond)
					continue
				}
			}

			if isDryrun {
				code = 999
				responseTime = time.Duration(0)
			} else {
				client := &http.Client{Timeout: time.Duration(20) * time.Second}
				httpreq, err := http.NewRequest(req.Method, req.URL, nil)
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
	r.wg.Done()
}

// Run all repeat requests
func (r *Repeater) Run(concurrency int, isDryrun bool, ignoreReqTime bool) {
	printer.Get().Run()
	go r.collectStats()

	for i := 0; concurrency > i; i++ {
		r.wg.Add(1)
		go r.runRequestWorker(isDryrun, ignoreReqTime)
	}

	for _, req := range *r.requests {
		r.buffer <- req
	}

	close(r.buffer)
	r.wg.Wait()
	close(resultBuffer)
	r.wgCollectStats.Wait()
	color.Blue("Repeat Completed!")
}
