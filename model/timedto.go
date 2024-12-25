package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = time.RFC3339

type TimeDTO time.Time

const (
	DurationYear  = DurationDay * 365
	DurationMonth = DurationYear / 12
	DurationDay   = time.Hour * 24
)

var durationNames = []string{"y", "m", "d"}

func (t TimeDTO) String() string {
	return time.Time(t).Format(TimeFormat)
}

func (t TimeDTO) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t *TimeDTO) UnmarshalText(text []byte) error {
	if d, err := parserAsDuration(string(text)); err == nil {
		nt := time.Now().Add(d)
		*t = TimeDTO(nt)
		return nil
	}
	// failed to parse as duration, try as absolute time
	tm, err := time.Parse(TimeFormat, string(text))
	if err != nil {
		return err
	}
	*t = TimeDTO(tm)
	return nil
}

func parserAsDuration(s string) (time.Duration, error) {
	if strings.EqualFold(s, "now") {
		return 0, nil
	}
	for _, d := range durationNames {
		if !strings.HasSuffix(s, d) {
			continue
		}
		i, err := strconv.Atoi(strings.TrimSuffix(s, d))
		if err != nil {
			return 0, err
		}
		size := time.Duration(i)
		switch d {
		case "y":
			return DurationYear * size, nil
		case "m":
			return DurationMonth * size, nil
		case "d":
			return DurationDay * size, nil
		default:
			return 0, fmt.Errorf("invalid time unit %q", d)
		}
	}
	return 0, fmt.Errorf("not a known duration %q", s)
}
