package request

import (
	"sort"
	"testing"
	"time"

	"github.com/yomon8/logrepeat/parser"
)

var testtime = time.Now()

var cases = []struct {
	method, host, port, path string
	reqtime                  time.Time
}{
	{"http", "localhost", "80", "index.html", testtime.Add(time.Hour * -48)},
	{"http", "localhost", "80", "index.html", testtime.Add(time.Hour * -1)},
	{"http", "localhost", "80", "path/q=あいうえお", testtime.Add(time.Hour * -1)},
	{"http", "localhost", "80", "index.html", testtime.Add(time.Minute * -3)},
	{"http", "localhost", "80", "index.html", testtime.Add(time.Hour * -43)},
}

func TestPrintRequest(t *testing.T) {
	var requests Requests = make([]*RequestEntry, 0)
	for _, c := range cases {
		entry := new(parser.Entry)
		entry.Method = c.method
		entry.DateTime = time.Now()
		entry.Path = c.path
		r := NewRequestEntry(
			c.host,
			c.port,
			entry)
		requests = append(requests, r)
	}

	t.Log("LIST---------")
	for _, r := range requests {
		t.Logf("%s\n", r)
	}

	t.Log("SORT---------")
	sort.Sort(requests)
	for _, r := range requests {
		t.Logf("%s\n", r)
	}

}
