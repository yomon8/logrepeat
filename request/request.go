package request

import (
	"fmt"
	"time"

	"github.com/yomon8/logrepeat/parser"
)

type RequestEntry struct {
	URL        string
	Method     string
	OriginTime time.Time
	RepeatTime time.Time
}

func NewRequestEntry(host, port string, entry *parser.Entry) *RequestEntry {
	r := new(RequestEntry)
	r.Method = entry.Method
	r.URL = fmt.Sprintf("%s://%s:%s/%s", entry.Protocol, host, port, entry.Path)
	r.OriginTime = entry.DateTime
	return r
}

func (r *RequestEntry) String() string {
	return fmt.Sprintf("[%s]%s", r.Method, r.URL)
}

var datetimePrintFormat = "2006-01-02 15:04:05MST"

func (r *RequestEntry) StringOriginTime() string {
	return fmt.Sprintf("%s", r.OriginTime.In(time.Local).Format(datetimePrintFormat))
}

func (r *RequestEntry) StringPlanTime() string {
	return fmt.Sprintf("%s", r.RepeatTime.In(time.Local).Format(datetimePrintFormat))
}
