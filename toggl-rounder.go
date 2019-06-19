package main

import (
	"flag"
	"github.com/ravbaker/toggl-rounder/internal/rounder"
	"os"
)

func main() {
	config := rounder.Config{
		Rounding: *flag.Bool("rounding", false, "Should it round last entry?"),
		DryRun:  *flag.Bool("dry", true, "Should it update toggl?"),
	}
	flag.Parse()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		println("Missing $API_KEY for Toggl!")
		os.Exit(-1)
	}
	rounder.RoundThisMonth(apiKey, &config);
}
