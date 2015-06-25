package time

import (
    "time"
)

type Timestamp int64

func NewTimestamp() Timestamp {
	return Timestamp(time.Now().Unix())
}

func TimestampWithTime(time time.Time) Timestamp {
	return Timestamp(time.Unix())
}

func (this *Timestamp) Time() time.Time {
	return time.Unix(int64(*this), 0)
}

func (this *Timestamp) Add(year, month, day int, hour, minute, second int64) Timestamp {
	var d int64 = second + minute*60 + hour*60*60
	t := this.Time()
	t = t.AddDate(year, month, day).Add(time.Duration(d) * time.Second)
	return TimestampWithTime(t)
}

func (this *Timestamp) AddDate(year, month, day int) Timestamp {
	return this.Add(year, month, day, 0, 0, 0)
}
