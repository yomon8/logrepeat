package reperter

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Result struct {
	ReturnCode   int
	RequestTime  time.Time
	ResponseTime time.Duration
}

func (r *Result) requestTimeString() string {
	var datetimePrintFormat = "2006-01-02T15:04:05MST"
	return fmt.Sprintf("%s", r.RequestTime.In(time.Local).Format(datetimePrintFormat))

}

type Results []*Result

func newResults(size int, statusCode int) Results {
	rs := make([]*Result, size)
	return rs
}

func (rs Results) GetResultCount() int {
	return len(rs)
}

func (rs Results) GetStatsString() string {
	var (
		status2XX, status3XX, status4XX, status5XX, other, errorStatus int
	)

	var totalResponseTime time.Duration
	for _, r := range rs {
		switch {
		case 200 <= r.ReturnCode && r.ReturnCode < 300:
			status2XX++
		case 300 <= r.ReturnCode && r.ReturnCode < 400:
			status3XX++
		case 400 <= r.ReturnCode && r.ReturnCode < 500:
			status4XX++
		case 500 <= r.ReturnCode && r.ReturnCode < 600:
			status5XX++
		case 1000 <= r.ReturnCode:
			errorStatus++
		default:
			other++
		}
		totalResponseTime += r.ResponseTime
	}
	avgResponseTime := totalResponseTime / time.Duration(rs.Len()) / time.Millisecond

	return fmt.Sprintf(
		"/%s /3xx:%4d /%s /%s /%s /%s /Avg %d msec",
		color.GreenString("2xx:%4d", status2XX),
		status3XX,
		color.YellowString("4xx:%4d", status4XX),
		color.RedString("5xx:%4d", status5XX),
		color.MagentaString("Oth:%4d", other),
		color.HiRedString("Err:%4d", errorStatus),
		avgResponseTime)
}

func (rs Results) add(r *Result) Results {
	rs = append(rs, r)
	return rs
}

func (rs Results) Len() int {
	return len(rs)
}

func (rs Results) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs Results) Less(i, j int) bool {
	return rs[i].RequestTime.Before(rs[j].RequestTime)
}

func newResult(requestTime time.Time, statusCode int, responseTime time.Duration) *Result {
	res := new(Result)
	res.RequestTime = requestTime
	res.ReturnCode = statusCode
	res.ResponseTime = responseTime
	return res
}
