package time

import (
	"time"
	"fmt"
	"gopkg.in/mgo.v2/bson"
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
	if this == 0 {
		return time.Unix(int64(this), 0)
	}
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
	return fmt.Sprintf("%.3f", this)
}

func (this Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(this.String()), nil
}

// 如果不用 mgo, 可以将下面的代码注释掉
func (this Timestamp) GetBSON() (interface{}, error) {
	if this.Time().IsZero() {
		return nil, nil
	}
	return this.Time(), nil
}

func (this *Timestamp) SetBSON(raw bson.Raw) (err error) {
	var tm time.Time
	if err = raw.Unmarshal(&tm); err == nil {
		*this = TimestampWithTime(tm)
		return nil
	}

	var tms float64
	if err = raw.Unmarshal(&tms); err == nil {
		*this = Timestamp(tms)
		return nil
	}

	return err
}

