package main

import (
	"flag"
	"os"

	"github.com/ravbaker/toggl-rounder/internal/rounder"
)

func main() {
	version := flag.Bool("version", false, "Print the version")
	apiKey := flag.String("api-key", os.Getenv("TOGGL_API_KEY"), "Toggl API KEY `secret-key`, can also be provided via $TOGGL_API_KEY environment variable")
	appConfig := rounder.Config{
		Rounding: *flag.Bool("rounding", false, "Enables rounding last entry up to full unit"),
		DryRun:   *flag.Bool("dry", true, "Unless set to false it doesn't update records in Toggl"),
		Debug:   *flag.Bool("debug", false, "Print debugging output of API calls"),
	}
	flag.Parse()

	if *version {
		rounder.PrintVersion()
		return
	}

	if *apiKey == "" {
		println("Missing value for -api-key","\t", flag.Lookup("api-key").Usage)
		os.Exit(-1)
	}
	rounder.RoundThisMonth(*apiKey, &appConfig)
}
