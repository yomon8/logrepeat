package parser

import (
	"errors"
	"time"
)

var (
	// ErrIgnored should be set, If parsed line should be ignored.
	ErrIgnored = errors.New("line ignored")
)

// LogParser is parser of original log file.
type LogParser interface {
	Parse(string) (*Entry, error)
}

// Entry is created from parsed and used by repeater
type Entry struct {
	// Method is url method ex. GET,POST...
	Method string
	// Protocol is protocol of request ex. http,https...
	Protocol string
	// DateTime is at first originl request time, repeater update it as repeat request timing
	DateTime time.Time
	// Path is url path
	Path string
}
