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
	appConfig = *NewConfig(false, false, "keep")
	fakeUpdater := &fakeTogglUpdater{}
	list := []toggl.TimeEntry{
		buildEntry(makeTime(time.June, 21, 9, 24, 0), 5*time.Hour+24*time.Minute+11*time.Second),
		buildEntry(makeTime(time.June, 20, 8, 4, 12), 6*time.Hour),
		buildEntry(makeTime(time.June, 19, 8, 4, 12), 5*time.Hour+4*time.Minute+42*time.Second),
	}

	updateEntries(list, fakeUpdater)

	if len(fakeUpdater.entries) != 4 {
		t.Errorf("Updates count mismatch. Expect 4, got: %d", len(fakeUpdater.entries))
	}

	entryFor19th := fakeUpdater.entries[0] // last from list
	entryFor20th := fakeUpdater.entries[1]
	entryFor21th := fakeUpdater.entries[2]
	entryFor21thUpdated := fakeUpdater.entries[3]

	if !matchTimeSpec(entryFor19th, makeTime(time.June, 19, 8, 0, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor19th entry updated: expected 8:00->13:00 with 5h duration while got: %s->%s", entryFor19th.Start.String(), entryFor19th.Stop.String())
	}

	if !matchTimeSpec(entryFor20th, makeTime(time.June, 20, 8, 0, 0), 6*time.Hour) {
		t.Errorf("Wrong entryFor20th entry updated: expected 8:00->14:00 with 6h duration while got: %s->%s", entryFor20th.Start.String(), entryFor20th.Stop.String())
	}
	if !matchTimeSpec(entryFor21th, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor21th update for 21th entry: expected 9:30->13:30 with 5h30 duration while got: %s->%s", entryFor21th.Start.String(), entryFor21th.Stop.String())
	}

	if !matchTimeSpec(entryFor21thUpdated, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour+28*time.Minute+53*time.Second) {
		t.Errorf("Wrong second update for 21th entry: expected 9:30->14:00 with 5h30 duration while got: %s->%s", entryFor21thUpdated.Start.String(), entryFor21thUpdated.Stop.String())
	}
}

func TestUpdateEntriesWithRoundUp(t *testing.T) {
	appConfig = *NewConfig(false, false, "round_up")
	fakeUpdater := &fakeTogglUpdater{}
	list := []toggl.TimeEntry{
		buildEntry(makeTime(time.June, 21, 9, 24, 0), 5*time.Hour+4*time.Minute+11*time.Second),
		buildEntry(makeTime(time.June, 20, 8, 4, 12), 6*time.Hour),
		buildEntry(makeTime(time.June, 19, 8, 4, 12), 5*time.Hour+4*time.Minute+42*time.Second),
	}

	updateEntries(list, fakeUpdater)

	if len(fakeUpdater.entries) != 4 {
		t.Errorf("Updates count mismatch. Expect 4, got: %d", len(fakeUpdater.entries))
	}

	entryFor19th := fakeUpdater.entries[0] // last from list
	entryFor20th := fakeUpdater.entries[1]
	entryFor21th := fakeUpdater.entries[2]
	entryFor21thUpdated := fakeUpdater.entries[3]

	if !matchTimeSpec(entryFor19th, makeTime(time.June, 19, 8, 0, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor19th entry updated: expected 8:00->13:00 with 5h duration while got: %s->%s", entryFor19th.Start.String(), entryFor19th.Stop.String())
	}

	if !matchTimeSpec(entryFor20th, makeTime(time.June, 20, 8, 0, 0), 6*time.Hour) {
		t.Errorf("Wrong entryFor20th entry updated: expected 8:00->14:00 with 6h duration while got: %s->%s", entryFor20th.Start.String(), entryFor20th.Stop.String())
	}
	if !matchTimeSpec(entryFor21th, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor21th update for 21th entry: expected 9:30->13:30 with 5h30 duration while got: %s->%s", entryFor21th.Start.String(), entryFor21th.Stop.String())
	}

	if !matchTimeSpec(entryFor21thUpdated, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour+30*time.Minute) {
		t.Errorf("Wrong second update for 21th entry: expected 9:30->14:00 with 5h30 duration while got: %s->%s", entryFor21thUpdated.Start.String(), entryFor21thUpdated.Stop.String())
	}
}

func TestUpdateEntriesWithRoundHalfBelow(t *testing.T) {
	appConfig = *NewConfig(false, false, "round_half")
	fakeUpdater := &fakeTogglUpdater{}
	list := []toggl.TimeEntry{
		buildEntry(makeTime(time.June, 21, 9, 24, 0), 5*time.Hour+4*time.Minute+11*time.Second),
		buildEntry(makeTime(time.June, 20, 8, 4, 12), 6*time.Hour),
		buildEntry(makeTime(time.June, 19, 8, 4, 12), 5*time.Hour+4*time.Minute+42*time.Second),
	}

	updateEntries(list, fakeUpdater)

	if len(fakeUpdater.entries) != 4 {
		t.Errorf("Updates count mismatch. Expect 4, got: %d", len(fakeUpdater.entries))
	}

	entryFor19th := fakeUpdater.entries[0] // last from list
	entryFor20th := fakeUpdater.entries[1]
	entryFor21th := fakeUpdater.entries[2]
	entryFor21thUpdated := fakeUpdater.entries[3]

	if !matchTimeSpec(entryFor19th, makeTime(time.June, 19, 8, 0, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor19th entry updated: expected 8:00->13:00 with 5h duration while got: %s->%s", entryFor19th.Start.String(), entryFor19th.Stop.String())
	}

	if !matchTimeSpec(entryFor20th, makeTime(time.June, 20, 8, 0, 0), 6*time.Hour) {
		t.Errorf("Wrong entryFor20th entry updated: expected 8:00->14:00 with 6h duration while got: %s->%s", entryFor20th.Start.String(), entryFor20th.Stop.String())
	}
	if !matchTimeSpec(entryFor21th, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor21th update for 21th entry: expected 9:30->13:30 with 5h30 duration while got: %s->%s", entryFor21th.Start.String(), entryFor21th.Stop.String())
	}

	if !matchTimeSpec(entryFor21thUpdated, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour) {
		t.Errorf("Wrong second update for 21th entry: expected 9:30->14:00 with 5h00 duration while got: %s->%s", entryFor21thUpdated.Start.String(), entryFor21thUpdated.Stop.String())
	}
}

func TestUpdateEntriesWithRoundHalfAbove(t *testing.T) {
	appConfig = *NewConfig(false, false, "round_half")
	fakeUpdater := &fakeTogglUpdater{}
	list := []toggl.TimeEntry{
		buildEntry(makeTime(time.June, 21, 9, 24, 0), 5*time.Hour+4*time.Minute+11*time.Second),
		buildEntry(makeTime(time.June, 20, 8, 4, 12), 6*time.Hour),
		buildEntry(makeTime(time.June, 19, 8, 4, 12), 5*time.Hour+10*time.Minute+50*time.Second),
	}

	updateEntries(list, fakeUpdater)

	if len(fakeUpdater.entries) != 4 {
		t.Errorf("Updates count mismatch. Expect 4, got: %d", len(fakeUpdater.entries))
	}

	entryFor19th := fakeUpdater.entries[0] // last from list
	entryFor20th := fakeUpdater.entries[1]
	entryFor21th := fakeUpdater.entries[2]
	entryFor21thUpdated := fakeUpdater.entries[3]

	if !matchTimeSpec(entryFor19th, makeTime(time.June, 19, 8, 0, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor19th entry updated: expected 8:00->13:00 with 5h duration while got: %s->%s", entryFor19th.Start.String(), entryFor19th.Stop.String())
	}

	if !matchTimeSpec(entryFor20th, makeTime(time.June, 20, 8, 0, 0), 6*time.Hour) {
		t.Errorf("Wrong entryFor20th entry updated: expected 8:00->14:00 with 6h duration while got: %s->%s", entryFor20th.Start.String(), entryFor20th.Stop.String())
	}
	if !matchTimeSpec(entryFor21th, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor21th update for 21th entry: expected 9:30->13:30 with 5h30 duration while got: %s->%s", entryFor21th.Start.String(), entryFor21th.Stop.String())
	}

	if !matchTimeSpec(entryFor21thUpdated, makeTime(time.June, 21, 9, 30, 0), 5*time.Hour+30*time.Minute) {
		t.Errorf("Wrong second update for 21th entry: expected 9:30->14:00 with 5h30 duration while got: %s->%s", entryFor21thUpdated.Start.String(), entryFor21thUpdated.Stop.String())
	}
}

func TestUpdateEntriesThatOverlapAfterAdjusting(t *testing.T) {
	appConfig = *NewConfig(false, false, "keep")
	fakeUpdater := &fakeTogglUpdater{}
	list := []toggl.TimeEntry{
		buildEntry(makeTime(time.July, 20, 12, 15, 16), 3*time.Hour+24*time.Minute+11*time.Second),
		buildEntry(makeTime(time.July, 20, 8, 0, 0), 4*time.Hour+6*time.Minute),
		buildEntry(makeTime(time.July, 19, 9, 4, 12), 5*time.Hour+24*time.Minute+42*time.Second),
	}

	updateEntries(list, fakeUpdater)

	if len(fakeUpdater.entries) != 4 {
		t.Errorf("Updates count mismatch. Expect 4, got: %d", len(fakeUpdater.entries))
	}

	entryFor19th := fakeUpdater.entries[0] // last from list
	firstEntryFor20th := fakeUpdater.entries[1]
	secondEntryFor20th := fakeUpdater.entries[2]
	secondEntryFor20thUpdated := fakeUpdater.entries[3]

	if !matchTimeSpec(entryFor19th, makeTime(time.July, 19, 9, 0, 0), 5*time.Hour) {
		t.Errorf("Wrong entryFor19th entry updated: expected 9:00->14:00 with 5h duration while got: %s->%s", entryFor19th.Start.String(), entryFor19th.Stop.String())
	}

	if !matchTimeSpec(firstEntryFor20th, makeTime(time.July, 20, 8, 0, 0), 4*time.Hour+30*time.Minute) {
		t.Errorf("Wrong firstEntryFor20th entry updated: expected 8:00->12:30 with 4h30m duration while got: %s->%s", firstEntryFor20th.Start.String(), firstEntryFor20th.Stop.String())
	}
	if !matchTimeSpec(secondEntryFor20th, makeTime(time.July, 20, 12, 30, 0), 3*time.Hour) {
		t.Errorf("Wrong secondEntryFor20th update for 21th entry: expected 12:30->16:00 with 3h duration while got: %s->%s", secondEntryFor20th.Start.String(), secondEntryFor20th.Stop.String())
	}

	if !matchTimeSpec(secondEntryFor20thUpdated, makeTime(time.July, 20, 12, 30, 0), 3*time.Hour+24*time.Minute+53*time.Second) {
		t.Errorf("Wrong second update for 20th 2nd entry: expected 12:30->16:00 with 3h30 duration while got: %s->%s", secondEntryFor20thUpdated.Start.String(), secondEntryFor20thUpdated.Stop.String())
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
