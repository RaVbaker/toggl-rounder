package rounder

import (
	"fmt"
	"time"

	"github.com/jason0x43/go-toggl"
)

type Config struct {
	Rounding, DryRun bool
}

const (
	Version = "0.0.2"
	Granularity = 30 * time.Minute
)

var (
	remainingSum time.Duration = 0
	appConfig    Config
)

func PrintVersion() {
	fmt.Println(Version)
}

func RoundThisMonth(apiKey string, config *Config) {
	appConfig = *config
	session := toggl.OpenSession(apiKey)

	now := time.Now()
	monthBegin := fmt.Sprintf("%d-%02d-01", now.Year(), now.Month())
	today := now.Format("2006-01-02")

	entries, err := fetchAccountEntries(session, monthBegin, today)

	if err != nil {
		println("ERR:", err)
		return
	}

	updateEntries(entries, session)
}

func fetchAccountEntries(session toggl.Session, since, until string) ([]toggl.TimeEntry, error) {
	var entries []toggl.TimeEntry
	account, _ := session.GetAccount()
	workspaceId := account.Data.Workspaces[0].ID

	currentPage := 1
	for {
		report, err := session.GetDetailedReport(workspaceId, since, until, currentPage)
		if err != nil {
			return nil, err
		}
		entriesCount := len(report.Data)

		for i := 0; i < entriesCount; i++ {
			detailedTimeEntry := report.Data[i]
			if detailedTimeEntry.Uid == account.Data.ID {
				entry := buildTimeEntryFromDetails(workspaceId, detailedTimeEntry)
				entries = append(entries, entry)
			}
		}
		if entriesCount == 0 || report.TotalCount < report.PerPage {
			break
		}
		currentPage++
	}
	return entries, nil
}

func buildTimeEntryFromDetails(workspaceId int, entry toggl.DetailedTimeEntry) toggl.TimeEntry {
	return toggl.TimeEntry{
		Wid:         workspaceId,
		ID:          entry.ID,
		Pid:         entry.Pid,
		Tid:         entry.Tid,
		Description: entry.Description,
		Tags:        entry.Tags,
		Start:       entry.Start,
		Stop:        entry.End,
		Duration:    entry.End.Unix() - entry.Start.Unix(),
		DurOnly:     false,
		Billable:    entry.Billable,
	}
}

func updateEntries(entries []toggl.TimeEntry, session toggl.Session) {
	var entry toggl.TimeEntry
	for i := len(entries) - 1; i >= 0; i-- { // iterate from oldest to latest
		entry = entries[i]
		roundedTime := distributeRemaining(entry)
		displayEntry(entry, roundedTime, remainingSum)
		updateEntry(session, &entry, seconds(roundedTime))
	}
	if remainingSum > (Granularity/2) || (remainingSum > 0 && appConfig.Rounding) {
		updateEntry(session, &entry, entry.Duration+seconds(Granularity))
	}
	println(fmt.Sprintf("%.2f", remainingSum.Minutes()))
}

func distributeRemaining(entry toggl.TimeEntry) time.Duration {
	roundedTime, remaining := missingTime(entry)

	remainingSum += remaining
	if remainingSum > Granularity { // distribute remaining
		roundedTime += Granularity
		remainingSum -= Granularity
	}
	return roundedTime
}

func missingTime(entry toggl.TimeEntry) (time.Duration, time.Duration) {
	exactDuration := entry.Stop.Unix() - entry.Start.Unix()
	remaining := exactDuration % seconds(Granularity)
	return time.Duration(exactDuration-remaining) * time.Second, time.Duration(remaining) * time.Second
}

func displayEntry(entry toggl.TimeEntry, roundedTime time.Duration, remaining time.Duration) {
	println(entry.Start.Format("2006-01-02"), entry.Description, ":", entry.Duration, fmt.Sprintf("%.1f %.2f", roundedTime.Hours(), remaining.Minutes()), entry.Start.Format("15:04:05"), "->", entry.Stop.Format("15:04:05"))
}

func seconds(duration time.Duration) int64 {
	return int64(duration.Seconds())
}

func updateEntry(session toggl.Session, entry *toggl.TimeEntry, newDuration int64) {
	entry.SetStartTime(entry.Start.Truncate(time.Hour), true)
	_ = entry.SetDuration(newDuration)
	println("UPDATING:", entry.Duration, entry.Start.Format("15:04:05"), "->", entry.Stop.Format("15:04:05"))
	if !appConfig.DryRun {
		_, err := session.UpdateTimeEntry(*entry)
		if err != nil {
			println("ERR:", entry.ID, err)
		}
	}
}
