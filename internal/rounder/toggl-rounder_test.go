package rounder

import (
	"testing"
	"time"
	
	"github.com/jason0x43/go-toggl"
)

type fakeTogglUpdater struct {
	entries []toggl.TimeEntry
}

func (session *fakeTogglUpdater) UpdateTimeEntry(timer toggl.TimeEntry) (toggl.TimeEntry, error) {
	session.entries = append(session.entries, timer)
	return timer, nil
}

func TestUpdateEntries(t *testing.T) {
	fakeUpdater := &fakeTogglUpdater{}
	
	list := []toggl.TimeEntry{
		buildEntry(makeTime(time.June, 20, 8, 4, 12), 6 * time.Hour),
		buildEntry(makeTime(time.June, 19, 8, 4, 12), 5 * time.Hour + 4*time.Minute + 42*time.Second),
	}
	updateEntries(list, fakeUpdater)
	if len(fakeUpdater.entries) != 2 {
		t.Errorf("Entries count mismatch. Expect 2, got: %d", len(fakeUpdater.entries))
	}
	
	first := fakeUpdater.entries[0] // last from list
	last := fakeUpdater.entries[1]
	
	if !matchTimeSpec(first, makeTime(time.June, 19, 8, 0, 0), 5*time.Hour) {
		t.Errorf("Wrong first entry updated: expected 8:00->13:00 with 5h duration while got: %s->%s", first.Start.String(), first.Stop.String())
	}
	
	if !matchTimeSpec(last, makeTime(time.June, 20, 8, 0, 0), 6*time.Hour) {
		t.Errorf("Wrong last entry updated: expected 8:00->14:00 with 6h duration while got: %s->%s", last.Start.String(), last.Stop.String())
	}
}

func buildEntry(startTime *time.Time, duration time.Duration) toggl.TimeEntry {
	endTime := startTime.Add(duration)
	return toggl.TimeEntry{
		Start: startTime,
		Stop:  &endTime,
	}
}

func makeTime(month time.Month, day, hour, min, sec int) *time.Time {
	obj := time.Date(2019, month, day, hour, min, sec, 0, time.UTC)
	return &obj
}

func matchTimeSpec(entry toggl.TimeEntry, expectedStart *time.Time, duration time.Duration) bool {
	expectedEnd := expectedStart.Add(duration)
	return entry.Start.Equal(*expectedStart) && entry.Stop.Equal(expectedEnd) && entry.Duration == seconds(duration)
}
