package main

import (
	"flag"
	"os"

	"github.com/ravbaker/toggl-rounder/internal/rounder"
)

func main() {
	version := flag.Bool("version", false, "Print the version")

	config := rounder.Config{
		Rounding: *flag.Bool("rounding", false, "Enables rounding last entry up to full unit"),
		DryRun:   *flag.Bool("dry", true, "Unless set to false it doesn't update records in Toggl"),
	}
	flag.Parse()

	if *version {
		rounder.PrintVersion()
		return
	}

	apiKey := os.Getenv("TOGGL_API_KEY")
	if apiKey == "" {
		println("Missing $TOGGL_API_KEY environment variable!")
		os.Exit(-1)
	}
	rounder.RoundThisMonth(apiKey, &config)
}
