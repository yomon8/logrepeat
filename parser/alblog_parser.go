package parser

import (
	"errors"
	"regexp"
	"time"
)

// ALBLogParser can parse AWS ALB Access Log
type ALBLogParser struct {
	format string
	regexp *regexp.Regexp
}

// NewALBLogParser create AlbLogParser instance
// set regexp refs below url
// http://docs.aws.amazon.com/en_us/elasticloadbalancing/latest/application/load-balancer-access-logs.html
func NewALBLogParser() *ALBLogParser {
	a := new(ALBLogParser)
	a.format = "2006-01-02T15:04:05.000000Z"
	a.regexp = regexp.MustCompile(
		`^.+? (\d{4}-\d{2}-\d{2}T\d{2}\:\d{2}\:\d{2}.\d{6}Z) ` + //Time
			`.+? .+? .+? .+? .+? .+? .+? .+? .+? .+? \"` +
			`(.+?) ` + //Method
			`([^:]+):\/{2,3}[0-9a-z\.\-:\[\]]+?:?[0-9]{0,5}?\/(|.+?) ` + //Protocol,URL
			`.+?\" .+? .+? .+? .+? \".+?\"$`)
	return a
}

// Parse log line
func (a *ALBLogParser) Parse(line string) (*Entry, error) {
	matches := a.regexp.FindStringSubmatch(line)
	if len(matches) < 4 {
		return nil, errors.New("parse error")
	}

	// ALBLog Parser only support GET requests
	method := matches[2]
	if method == "POST" {
		return nil, ErrNoSupport
	}

	prot := matches[3]
	if prot != "http" && prot != "https" {
		return nil, ErrNoSupport
	}

	dt, err := time.Parse(a.format, matches[1])
	if err != nil {
		return nil, errors.New("time parse error")
	}
	return &Entry{
		DateTime: dt,
		Method:   matches[2],
		Protocol: prot,
		Path:     matches[4],
	}, nil
}
