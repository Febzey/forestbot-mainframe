package utils

import (
	"strconv"
	"time"
)

var Epoch = time.Unix(0, 0)

func MsToTime(ms string) (time.Time, error) {
	ts, err := strconv.Atoi(ms)
	if err != nil {
		return time.Time{}, err
	}

	return Epoch.Add(time.Duration(ts) * time.Millisecond), nil
}
