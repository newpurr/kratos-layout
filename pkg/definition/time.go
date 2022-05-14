package definition

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type (
	LocalTime struct {
		time.Time
	}
)

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t LocalTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// UnmarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t *LocalTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	tm, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = LocalTime{tm}
	return nil
}

// Value insert timestamp into mysql need this function.
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
