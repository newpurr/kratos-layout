package timeutils

import "time"

type TimeInterval [2]time.Time

func (t TimeInterval) GetBeg() time.Time {
	return t[0]
}

func (t TimeInterval) GetEnd() time.Time {
	return t[1]
}

type TimeIntervalSearcher interface {
	GetBeg() string
	GetEnd() string
}

func ParseDatetimeInterval(t TimeIntervalSearcher) *TimeInterval {
	if t.GetBeg() == "" || t.GetEnd() == "" {
		return nil
	}
	r1, err1 := time.Parse("2006-01-02 15:04:05", t.GetBeg())
	r2, err2 := time.Parse("2006-01-02 15:04:05", t.GetEnd())
	if err1 != nil && err2 != nil {
		return (*TimeInterval)(&([2]time.Time{
			r1, r2,
		}))
	}
	return nil
}

func ParseDateInterval(t TimeIntervalSearcher) [2]time.Time {
	if t.GetBeg() == "" || t.GetEnd() == "" {
		return [2]time.Time{}
	}
	r1, err1 := time.Parse("2006-01-02", t.GetBeg())
	r2, err2 := time.Parse("2006-01-02", t.GetEnd())
	if err1 != nil && err2 != nil {
		return [2]time.Time{
			r1, r2,
		}
	}
	return [2]time.Time{}
}
