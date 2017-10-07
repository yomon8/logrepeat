package request

import (
	"sort"
	"time"
)

type Requests []*RequestEntry

func (reqs Requests) UpdateRepeatTime(differenceFromOriginalTime time.Duration) {
	for _, request := range reqs {
		request.RepeatTime = request.OriginTime.Add(differenceFromOriginalTime)
	}
}

func (reqs Requests) Len() int {
	return len(reqs)
}

func (reqs Requests) Swap(i, j int) {
	reqs[i], reqs[j] = reqs[j], reqs[i]
}

func (reqs Requests) Less(i, j int) bool {
	return reqs[i].OriginTime.Before(reqs[j].OriginTime)
}

func (reqs Requests) SortByOriginalTime() {
	sort.Sort(reqs)
}
