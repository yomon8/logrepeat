package parser

import (
	"errors"
	"time"
)

var (
	ErrIgnored = errors.New("line ignored")
)

type LogParser interface {
	Parse(string) (*Entry, error)
}

type Entry struct {
	Method   string
	DateTime time.Time
	Path     string
}
