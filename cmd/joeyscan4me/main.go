package main

import (
	"os"

	"github.com/henrique-gomesz/joeyscan4me/pkg/logging"
	"github.com/henrique-gomesz/joeyscan4me/pkg/runner"
)

var Version string = "1.1.2"

func main() {
	logging.PrintBanner(Version)

	opt := runner.ParseOptions()

	if err := runner.StartScan(opt); err != nil {
		logging.LogError("Scan failed", err)
		os.Exit(1)
	}
}
