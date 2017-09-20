package parser

import (
	"errors"
	"regexp"
	"time"
)

type AlbLogParser struct {
	format string
	regexp *regexp.Regexp
}

func NewAlbLogParser() *AlbLogParser {
	a := new(AlbLogParser)
	a.format = "2006-01-02T15:04:05.000000Z"
	a.regexp = regexp.MustCompile(
		`^(.+?) ` + //Method
			`(\d{4}-\d{2}-\d{2}T\d{2}\:\d{2}\:\d{2}.\d{6}Z) ` + //Time
			`.+? .+? .+? .+? .+? .+? .+? .+? .+? .+? \".+? ` +
			`[^:]+:\/{2,3}[0-9a-z\.\-:\[\]]+?:?[0-9]{0,5}?\/(|.+?) ` + //URL
			`.+?\" .+? .+? .+? .+? \".+?\"$`)
	return a
}

func (a *AlbLogParser) Parse(line string) (*Entry, error) {
	matches := a.regexp.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil, errors.New("parse error")
	}
	entry := new(Entry)
	entry.Method = matches[1]
	var err error
	entry.DateTime, err = time.Parse(a.format, matches[2])
	if err != nil {
		return nil, errors.New("time parse error")
	}
	entry.Path = matches[3]
	return entry, nil
}
