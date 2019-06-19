package main

import (
	"flag"
	"os"

	"github.com/ravbaker/toggl-rounder/internal/rounder"
)

func main() {
	config := rounder.Config{
		Rounding: *flag.Bool("rounding", false, "Should it round last entry?"),
		DryRun:   *flag.Bool("dry", true, "Should it update toggl?"),
	}
	flag.Parse()

	apiKey := os.Getenv("TOGGL_API_KEY")
	if apiKey == "" {
		println("Missing $TOGGL_API_KEY environment variable!")
		os.Exit(-1)
	}
	rounder.RoundThisMonth(apiKey, &config)
}
