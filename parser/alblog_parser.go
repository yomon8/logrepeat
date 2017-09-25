package parser

import (
	"errors"
	"regexp"
	"time"
)

// AlbLogParser can parse AWS ALB Access Log
type AlbLogParser struct {
	format string
	regexp *regexp.Regexp
}

// NewAlbLogParser create AlbLogParser instance
// set regexp refs below url
// http://docs.aws.amazon.com/en_us/elasticloadbalancing/latest/application/load-balancer-access-logs.html
func NewAlbLogParser() *AlbLogParser {
	a := new(AlbLogParser)
	a.format = "2006-01-02T15:04:05.000000Z"
	a.regexp = regexp.MustCompile(
		`^(.+?) ` + //Protocol
			`(\d{4}-\d{2}-\d{2}T\d{2}\:\d{2}\:\d{2}.\d{6}Z) ` + //Time
			`.+? .+? .+? .+? .+? .+? .+? .+? .+? .+? \"` +
			`(.+?) ` + //Method
			`[^:]+:\/{2,3}[0-9a-z\.\-:\[\]]+?:?[0-9]{0,5}?\/(|.+?) ` + //URL
			`.+?\" .+? .+? .+? .+? \".+?\"$`)
	return a
}

// Parse log line
func (a *AlbLogParser) Parse(line string) (*Entry, error) {
	matches := a.regexp.FindStringSubmatch(line)
	if len(matches) < 4 {
		return nil, errors.New("parse error")
	}
	entry := new(Entry)
	entry.Protocol = matches[1]
	var err error
	entry.DateTime, err = time.Parse(a.format, matches[2])
	if err != nil {
		return nil, errors.New("time parse error")
	}
	entry.Method = matches[3]
	entry.Path = matches[4]
	return entry, nil
}
