package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/yomon8/logrepeat/parser"
	"github.com/yomon8/logrepeat/repeater"
	"github.com/yomon8/logrepeat/request"
)

const (
	// default values
	defaultSampleCount  = 5
	defaultHost         = "localhost"
	defaultPort         = "80"
	defaultConcurrency  = 10
	defaultAfterSeconds = 5
)

var (
	// args
	host              string
	port              string
	file              string
	samplecount       int
	concurrency       int
	afterSeconds      int
	ignoreRequestTime bool
	isForceMode       bool
	isDryrun          bool
	isHelp            bool
	isVersion         bool

	readreqs     request.Requests
	ignoredLine  int
	parseErrLine int
	newest       *request.Request
	oldest       *request.Request
)

func parseArgs() {
	flag.StringVar(&host, "h", defaultHost, "Repert target hostname")
	flag.StringVar(&port, "p", defaultPort, "Repert target port number")
	flag.StringVar(&file, "f", "", "AWS ALB log file path")
	flag.IntVar(&samplecount, "s", defaultSampleCount, "A number of request samples at repeat plan")
	flag.IntVar(&concurrency, "c", defaultConcurrency, "Concurrency of requesters")
	flag.IntVar(&afterSeconds, "after-seconds", defaultAfterSeconds, "Repeat start after seconds")
	flag.BoolVar(&ignoreRequestTime, "ignore-req-time", false, "Ignore request time, send request in order of rows.")
	flag.BoolVar(&isForceMode, "force", false, "Force mode,Show no prompt")
	flag.BoolVar(&isDryrun, "dryrun", false, "dryrun")
	flag.BoolVar(&isHelp, "help", false, "Show help message")
	flag.BoolVar(&isVersion, "v", false, "Show version info")
	flag.Parse()
	if len(os.Args) == 1 || isHelp {
		flag.Usage()
		os.Exit(-1)
	}
	if isVersion {
		fmt.Println("version:", version)
		os.Exit(0)
	}
}

func main() {
	parseArgs()

	// open log file
	var sc *bufio.Scanner
	if file != "" {
		fp, err := os.Open(file)
		if err != nil {
			log.Fatalf("File open error:%#v\n", err)
		}
		sc = bufio.NewScanner(fp)
		defer fp.Close()
	} else {
		sc = bufio.NewScanner(os.Stdin)
	}

	// parse log
	readreqs = make([]*request.Request, 0)
	p := parser.NewALBLogParser()
	for i := 0; ; i++ {
		if sc.Scan() {
			s := sc.Text()
			entry, err := p.Parse(s)
			switch {
			case err == parser.ErrIgnored:
				ignoredLine++
				continue
			case err != nil:
				parseErrLine++
				continue
			}
			req := request.NewRequest(host, port, entry)
			readreqs = append(readreqs, req)
		} else {
			break
		}
	}
	if len(readreqs) == 0 {
		log.Fatalf("no valid entries found")
	}

	// update time of requests according to current time
	readreqs.SortByOriginalTime()
	oldest = readreqs[0]
	newest = readreqs[len(readreqs)-1]
	difftime := time.Now().Add(time.Duration(afterSeconds) * time.Second).Sub(oldest.OriginTime)
	readreqs.UpdateRepeatTime(difftime)

	if !isForceMode {
		// print repeat plan
		printStartMessage()
		// wait for user prompt
		var key string
		var ok bool
		for !ok {
			fmt.Print(color.MagentaString("Start/Cancel>"))
			fmt.Scanf("%s", &key)
			switch key {
			case "S", "s", "Start", "start":
				ok = true
			case "C", "c", "Cancel", "cancel":
				fmt.Println("canceled.")
				os.Exit(1)
			default:
				continue
			}
		}
	}

	// start repeater
	repeater := reperter.NewRepeater(&readreqs)
	repeater.Run(concurrency, isDryrun, ignoreRequestTime)
}
