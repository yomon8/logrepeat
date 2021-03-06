package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yomon8/logrepeat/parser"
	"github.com/yomon8/logrepeat/repeater"
	"github.com/yomon8/logrepeat/request"
)

const (
	// default values
	defaultSampleCount  = 5
	defaultHost         = "localhost"
	defaultHttpPort     = "80"
	defaultHttpsPort    = "443"
	defaultConcurrency  = 10
	defaultAfterSeconds = 5
)

var (
	// args
	host              string
	httpPort          string
	httpsPort         string
	file              string
	samplecount       int
	concurrency       int
	afterSeconds      int
	ignoreRequestTime bool
	forceHttp         bool
	isForceMode       bool
	isDryrun          bool
	isHelp            bool
	isVersion         bool

	readreqs        request.Requests
	ignoredLine     int
	nonSuportedLine int
	parseErrLine    int
	newest          *request.RequestEntry
	oldest          *request.RequestEntry
)

func parseArgs() {
	flag.StringVar(&host, "h", defaultHost, "Repert target hostname")
	flag.StringVar(&httpPort, "http-port", defaultHttpPort, "Repert target http port number")
	flag.StringVar(&httpsPort, "https-port", defaultHttpsPort, "Repert target https port number")
	flag.StringVar(&file, "f", "", "AWS ALB log file path")
	flag.IntVar(&samplecount, "s", defaultSampleCount, "A number of request samples at repeat plan")
	flag.IntVar(&concurrency, "c", defaultConcurrency, "Concurrency of requesters")
	flag.IntVar(&afterSeconds, "start-after-secs", defaultAfterSeconds, "Repeat start after seconds")
	flag.BoolVar(&ignoreRequestTime, "ignore-timestamp", false, "Ignore request timestamp, simply send request in order of rows.")
	flag.BoolVar(&forceHttp, "force-http", false, "Send all of requests as http(not https)")
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
		fmt.Println("version: ", version)
		os.Exit(0)
	}
}

func validateArgs() {
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
	readreqs = make([]*request.RequestEntry, 0)
	p := parser.NewALBLogParser()
	for i := 0; ; i++ {
		if sc.Scan() {
			s := sc.Text()
			entry, err := p.Parse(s)
			switch {
			case err == parser.ErrIgnored:
				ignoredLine++
				continue
			case err == parser.ErrNoSupport:
				nonSuportedLine++
				continue
			case err != nil:
				parseErrLine++
				continue
			}

			var port string
			switch {
			case entry.Protocol == "http":
				port = httpPort
			case entry.Protocol == "https" && forceHttp:
				entry.Protocol = "http"
				port = httpPort
			case entry.Protocol == "https" && !forceHttp:
				port = httpsPort
			default:
				nonSuportedLine++
				continue
			}

			req := request.NewRequestEntry(host, port, entry)
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

	// skip prompt in force mode
	if !isForceMode {
		printStartMessage()
		waitPrompt()
	}

	// start repeater
	repeater := reperter.NewRepeater(&readreqs)
	repeater.Run(concurrency, isDryrun, ignoreRequestTime)
}
