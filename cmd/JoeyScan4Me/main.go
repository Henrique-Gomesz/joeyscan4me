package main

import (
	"os"

	"github.com/Henrique-Gomesz/JoeyScan4Me/pkg/logging"
	"github.com/Henrique-Gomesz/JoeyScan4Me/pkg/runner"
)

var Version string = "1.1.0"

func main() {
	logging.PrintBanner(Version)

	opt := runner.ParseOptions()

	if err := runner.StartScan(opt); err != nil {
		logging.LogError("Scan failed", err)
		os.Exit(1)
	}
}
