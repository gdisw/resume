package testutil

import "time"

func ParseDate(in string) time.Time {
	t, _ := time.Parse("2006-01-02", in)
	return t
}

func ParseDateTime(in string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", in)
	if err != nil {
		t, _ = time.Parse("2006-01-02 15:04", in)
	}

	return t
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func StartOfWeek(t time.Time) time.Time {
	runner := t
	for runner.Weekday() != time.Monday {
		runner = runner.AddDate(0, 0, -1)
	}

	return StartOfDay(runner)
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0,
		t.Location(),
	)
}
