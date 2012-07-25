package stat

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

func TestSimplejson(t *testing.T) {
	stat, _ := New("test")

	for i := 0; i < 20000; i++ {
		stat.Store(uint64(i), time.Now())
	}

	ret := stat.Calc()
	assert.Equal(t, ret.Average, uint64(17499))
	assert.Equal(t, ret.NinetyFivePercent, uint64(19751))
	assert.Equal(t, ret.NinetyNinePercent, uint64(19951))
	assert.Equal(t, ret.HundredPercent, uint64(19999))
	assert.Equal(t, ret.Count, uint64(20000))
}
