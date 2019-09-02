package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func FloatFromString(raw interface{}) (float64, error) {
	str, ok := raw.(string)
	if !ok {
		return 0, errors.New(fmt.Sprintf("unable to parse, value not string: %T", raw))
	}
	flt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("unable to parse as float: %s", str))
	}
	return flt, nil
}

func IntFromString(raw interface{}) (int, error) {
	str, ok := raw.(string)
	if !ok {
		return 0, errors.New(fmt.Sprintf("unable to parse, value not string: %T", raw))
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("unable to parse as int: %s", str))
	}
	return n, nil
}

func TimeFromUnixTimestampString(raw interface{}) (time.Time, error) {
	str, ok := raw.(string)
	if !ok {
		return time.Time{}, errors.New(fmt.Sprintf("unable to parse, value not string"))
	}
	ts, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, fmt.Sprintf("unable to parse as int: %s", str))
	}
	return time.Unix(0, ts*int64(time.Millisecond)), nil
}

func TimeFromUnixTimestampFloat(raw interface{}) (time.Time, error) {
	ts, ok := raw.(float64)
	if !ok {
		return time.Time{}, errors.New(fmt.Sprintf("unable to parse, value not int64: %T", raw))
	}
	return time.Unix(0, int64(ts)*int64(time.Millisecond)), nil
}

func UnixMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func RecvWindow(d time.Duration) int64 {
	return int64(d) / int64(time.Millisecond)
}

func Decimal(value float64, v string) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+v+"f", value), 64)
	return value
}
