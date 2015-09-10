package time

import (
	"time"
	"fmt"
)

type Timestamp float64

func NewTimestamp() Timestamp {
	return TimestampWithTime(time.Now())
}

func TimestampWithTime(time time.Time) Timestamp {
	var t = float64(time.Unix()) + (float64(time.Nanosecond()) / 1e9)
	return Timestamp(t)
}

func (this Timestamp) Time() time.Time {
	return time.Unix(int64(this), int64(this * 1e9) % int64(this))
}

func (this Timestamp) Add(year, month, day int, hour, minute, second int64) Timestamp {
	var d int64 = second + minute*60 + hour*60*60
	t := this.Time()
	t = t.AddDate(year, month, day).Add(time.Duration(d) * time.Second)
	return TimestampWithTime(t)
}

func (this Timestamp) AddDate(year, month, day int) Timestamp {
	return this.Add(year, month, day, 0, 0, 0)
}

func (this Timestamp) String() string {
	return fmt.Sprintf("%f", this)
}

func (this Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(this.String()), nil
}