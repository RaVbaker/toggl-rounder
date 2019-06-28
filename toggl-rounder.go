package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gookit/color"

	"github.com/ravbaker/toggl-rounder/internal/rounder"
)

var version, colors, dryRun, debugMode *bool
var apiKey, remainingStrategy *string

func main() {
	parseArgs()
	appConfig := rounder.NewConfig(*dryRun, *debugMode, *remainingStrategy)

	if !*colors {
		color.Disable()
	}

	if *version {
		rounder.PrintVersion()
		return
	}

	rounder.RoundThisMonth(*apiKey, appConfig)
}

func parseArgs() {
	version = flag.Bool("version", false, "Print the version")
	colors = flag.Bool("colors", true, "Display colorful output in Terminal")
	apiKey = flag.String("api-key", os.Getenv("TOGGL_API_KEY"), "Toggl API KEY `secret-key`, can also be provided via $TOGGL_API_KEY environment variable")
	dryRun = flag.Bool("dry", true, "Unless set to false it doesn't update records in Toggl")
	remainingStrategy = flag.String("remaining", "keep", fmt.Sprintf("Decides on what to do with remaining time. Possible options: %q", rounder.AllowedRemainingStrategies))
	debugMode = flag.Bool("debug", false, "Print debugging output of API calls")
	flag.Parse()
	
	if *version {
		return
	}

	if !rounder.IsAllowedRemainingStrategy(*remainingStrategy) {
		fmt.Printf("Not allowed -remaining value: '%s'. Allowed: %q\n", *remainingStrategy, rounder.AllowedRemainingStrategies)
		os.Exit(-1)
	}
	if *apiKey == "" {
		fmt.Println("Missing value for -api-key", "\t", flag.Lookup("api-key").Usage)
		os.Exit(-1)
	}
}
