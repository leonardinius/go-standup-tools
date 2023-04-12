package main_test

import (
	"testing"
	"time"

	main "github.com/leonardinius/go-standup-tools/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNaturalDateParse(t *testing.T) {
	datetimeParse := func(input string) time.Time {
		if value, err := time.Parse(time.DateOnly, input); err == nil {
			return value
		}
		if value, err := time.Parse(time.DateTime, input); err == nil {
			return value
		}
		t.Errorf("unable to parse %s as date time", input)
		return time.Time{}
	}

	today := datetimeParse("2023-11-05")

	suite := []struct {
		name         string
		input        string
		expectedTime time.Time
		err          string
	}{
		{"empty string", "", time.Time{}, "parse error"},
		{"now is now", "now", today, ""},
		{"relative -1", "yesterday", datetimeParse("2023-11-04"), ""},
		{"relative -1 ago", "1 day ago", datetimeParse("2023-11-04"), ""},
		{"relative -2 ago", "2 days ago", datetimeParse("2023-11-03"), ""},
		{"relative future +2 ago", "in 2 days", datetimeParse("2023-11-07"), ""},
		{"relative month", "February", datetimeParse("2023-02-05"), ""},
		{"relative month 1st", "February 1st", datetimeParse("2023-02-01"), ""},
		{"relative month nTh", "February 13th", datetimeParse("2023-02-13"), ""},
		{"date", "2023-02-13", datetimeParse("2023-02-13"), ""},
		{"date time", "2023-02-13 01:02:03", datetimeParse("2023-02-13 01:02:03"), ""},
	}

	for _, test := range suite {
		t.Run(test.name, func(t *testing.T) {
			parsed, err := main.ParseNaturalDate(test.input, today)
			if test.err == "" {
				assert.Equal(t, test.expectedTime, parsed)
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.err)
			}
		})
	}
}
