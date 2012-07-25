package stat

import (
	"errors"
	"math"
	"sort"
	"time"
)

const statWindowCount = 5000
const statWindowDuration = 5 * time.Minute

type Point struct {
	Val uint64
	Ts  time.Time
}

type Stat struct {
	Name  string
	Count uint64
	Index uint16
	Data  []Point
}

type Return struct {
	Count             uint64
	Average           uint64
	HundredPercent    uint64
	NinetyNinePercent uint64
	NinetyFivePercent uint64
}

type Uint64Slice []uint64

func (s Uint64Slice) Len() int {
	return len(s)
}

func (s Uint64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Uint64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

var stats = make(map[string]*Stat)

func New(name string) (*Stat, error) {
	_, ok := stats[name]
	if ok {
		return nil, errors.New("DUPLICATE_STAT")
	}

	stat := &Stat{
		Name: name,
		Data: make([]Point, statWindowCount),
	}

	stats[name] = stat

	return stat, nil
}

func (s *Stat) Store(val uint64, ts time.Time) {
	point := &s.Data[s.Index]
	point.Val = val
	point.Ts = ts

	s.Count++
	s.Index++

	if s.Index >= statWindowCount {
		s.Index = 0
	}
}

func (s *Stat) Calc() *Return {
	var startIndex uint16

	now := time.Now()
	values := make(Uint64Slice, 0, statWindowCount)

	if s.Count <= statWindowCount {
		startIndex = 0
	} else {
		startIndex = s.Index
	}

	total := uint64(0)
	for c := uint16(0); (c < statWindowCount) || (s.Count <= statWindowCount && c < s.Index); c++ {
		point := s.Data[(startIndex+c)%statWindowCount]
		duration := now.Sub(point.Ts)
		if duration < statWindowDuration {
			values = append(values, point.Val)
			total += point.Val
		}
	}

	ret := &Return{
		Count: s.Count,
	}

	numValues := len(values)
	if numValues > 0 {
		ret.Average = total / uint64(numValues)

		sort.Sort(values)
		ret.HundredPercent = percentile(100.0, values, numValues)
		ret.NinetyNinePercent = percentile(99.0, values, numValues)
		ret.NinetyFivePercent = percentile(95.0, values, numValues)
	}

	return ret
}

func StoreValue(name string, val uint64, ts time.Time) error {
	stat, ok := stats[name]
	if !ok {
		return errors.New("INVALID_STAT")
	}
	stat.Store(val, ts)
	return nil
}

func StoreDuration(name string, startTs time.Time) error {
	endTs := time.Now()
	duration := endTs.Sub(startTs)
	return StoreValue(name, uint64(duration), endTs)
}

func percentile(perc float64, arr Uint64Slice, length int) uint64 {
	indexOfPerc := int(math.Ceil(((perc / 100.0) * float64(length)) + 0.5))
	if indexOfPerc >= length {
		indexOfPerc = length - 1
	}
	return arr[indexOfPerc]
}
