package rounder

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/jason0x43/go-toggl"
)

type Config struct {
	DryRun, Debug     bool
	RemainingStrategy string
}

func NewConfig(dryRun bool, debug bool, remainingStrategy string) *Config {
	remainingSum = 0
	lastEntryEnd = time.Time{}
	return &Config{DryRun: dryRun, Debug: debug, RemainingStrategy: remainingStrategy}
}

type togglUpdater interface {
	UpdateTimeEntry(timer toggl.TimeEntry) (toggl.TimeEntry, error)
}

const (
	Version     = "0.1.2"
	Granularity = 30 * time.Minute
)

var (
	AllowedRemainingStrategies = [...]string{"keep", "round_half", "round_up"}

	remainingSum time.Duration = 0
	lastEntryEnd time.Time
	appConfig    Config
)

func PrintVersion() {
	fmt.Println(Version)
}

func IsAllowedRemainingStrategy(candidate string) bool {
	for _, n := range AllowedRemainingStrategies {
		if candidate == n {
			return true
		}
	}
	return false
}

func RoundThisMonth(apiKey string, config *Config) {
	appConfig = *config
	if !appConfig.Debug {
		toggl.DisableLog()
	}
	session := toggl.OpenSession(apiKey)

	now := time.Now()
	monthBegin := fmt.Sprintf("%d-%02d-01", now.Year(), now.Month())
	today := now.Format("2006-01-02")

	entries, err := fetchAccountEntries(session, monthBegin, today)

	if err != nil {
		fmt.Println("ERR:", err)
		return
	}

	updateEntries(entries, &session)
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
		Duration:    entry.Duration / 1000, // this is rounded duration since go-toggl fetches only such
		DurOnly:     false,
		Billable:    entry.Billable,
	}
}

func updateEntries(entries []toggl.TimeEntry, session togglUpdater) {
	remainingSum = 0
	var entry toggl.TimeEntry
	for i := len(entries) - 1; i >= 0; i-- { // iterate from oldest to latest
		entry = entries[i]
		roundedTime := distributeRemaining(entry)
		updateEntry(session, &entry, seconds(roundedTime))
	}
	extraDuration := lastEntryRemainingDuration()
	updateEntry(session, &entry, entry.Duration+seconds(extraDuration))
	color.Success.Printf("=> Remaining time(strategy: %s): %s, recorded: %s\n", appConfig.RemainingStrategy, remainingSum, extraDuration)
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
	exactDuration := actualDuration(&entry)
	remaining := exactDuration % seconds(Granularity)
	return secondsAsDuration(exactDuration - remaining), secondsAsDuration(remaining)
}

func displayEntry(entry toggl.TimeEntry, roundedTime time.Duration) {
	fmt.Println(
		fmt.Sprintf("ENTRY<%s>", color.Question.Render(entry.ID)),
		color.Danger.Render(entry.Start.Format("2006-01-02")),
		"existing(rounded):", secondsAsDuration(entry.Duration).String(),
		"existing(actual):", secondsAsDuration(actualDuration(&entry)).String(),
		"expected:", roundedTime.String(),
		"remaining:", remainingSum.String(),
		entry.Start.Format("15:04:05"), "->", entry.Stop.Format("15:04:05"),
		"\n",
		color.Comment.Render(entry.Description),
	)
}

func seconds(duration time.Duration) int64 {
	return int64(duration.Seconds())
}

func secondsAsDuration(seconds int64) time.Duration {
	return time.Duration(seconds) * time.Second
}

func updateEntry(session togglUpdater, entry *toggl.TimeEntry, newDuration int64) {
	newStartTime := entry.Start.Round(Granularity)
	// fix starting point to not overlap last entry
	if newStartTime.Before(lastEntryEnd) && !entry.Stop.Equal(lastEntryEnd) {
		newStartTime = lastEntryEnd
	}
	// nothing changed so let's skip update
	if newStartTime.Equal(*entry.Start) && actualDuration(entry) == newDuration {
		return
	}
	displayEntry(*entry, secondsAsDuration(newDuration))
	entry.SetStartTime(newStartTime, true)
	_ = entry.SetDuration(newDuration)
	lastEntryEnd = *entry.Stop
	fmt.Println("UPDATING:", secondsAsDuration(entry.Duration).String(), entry.Start.Format("15:04:05"), "->", entry.Stop.Format("15:04:05"))
	if !appConfig.DryRun {
		_, err := session.UpdateTimeEntry(*entry)
		if err != nil {
			fmt.Println("ERR:", entry.ID, err)
		}
	}
}

func actualDuration(entry *toggl.TimeEntry) int64 {
	return entry.Stop.Unix() - entry.Start.Unix()
}

func lastEntryRemainingDuration() time.Duration {
	var extraDuration time.Duration
	switch appConfig.RemainingStrategy {
	case "round_half":
		if remainingSum > (Granularity / 2) {
			extraDuration = Granularity
		}
	case "round_up":
		if remainingSum > 0 {
			extraDuration = Granularity
		}
	case "keep":
		extraDuration = remainingSum
	default:
		fmt.Printf("Unknown remaining strategy: '%s'", appConfig.RemainingStrategy)
		os.Exit(-2)
	}
	return extraDuration
}
