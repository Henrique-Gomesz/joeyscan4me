package runner

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/henrique-gomesz/joeyscan4me/pkg/logging"

	"github.com/sensepost/gowitness/pkg/log"
	goWitnessRunner "github.com/sensepost/gowitness/pkg/runner"
	driver "github.com/sensepost/gowitness/pkg/runner/drivers"
	"github.com/sensepost/gowitness/pkg/writers"
	"github.com/sensepost/gowitness/web"
)

// thx Claude Sonnet
func RunGowitness(opt *Options) error {
	httpxOutputFile := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), HttpxOutputFile)
	screenshotsPath := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), "screenshots")
	dbPath := filepath.Join(screenshotsPath, "gowitness.sqlite3")

	if opt.Resume && FileNonEmpty(dbPath) {
		logging.LogInfo("Skipping gowitness — database already exists: " + dbPath)
		return nil
	}

	if _, err := os.Stat(httpxOutputFile); os.IsNotExist(err) {
		return fmt.Errorf("httpx output file not found at %s", httpxOutputFile)
	}

	urls, err := ReadFileLines(httpxOutputFile)
	if err != nil {
		return fmt.Errorf("failed to read URLs from httpx output: %w", err)
	}
	urls = NormalizeAndDedupeLines(urls)

	if len(urls) == 0 {
		logging.LogInfo("No URLs found to screenshot")
		return nil
	}

	logging.LogInfo("Running gowitness to capture screenshots from up domains")

	// Ensure screenshots directory exists
	if err := os.MkdirAll(screenshotsPath, 0755); err != nil {
		return fmt.Errorf("failed to create screenshots directory: %w", err)
	}

	// Configure gowitness options
	options := goWitnessRunner.NewDefaultOptions()
	options.Scan.ScreenshotPath = screenshotsPath
	options.Scan.ScreenshotFormat = "jpeg"
	options.Scan.ScreenshotJpegQuality = 80
	options.Scan.Threads = 50
	options.Scan.Timeout = 30
	options.Scan.Delay = 3
	options.Logging.LogScanErrors = true

	// Create slog logger
	logger := slog.New(log.Logger)

	// Initialize chromedp driver
	gwDriver, err := driver.NewChromedp(logger, *options)
	if err != nil {
		return fmt.Errorf("failed to create gowitness driver: %w", err)
	}
	defer gwDriver.Close()

	// Create database writer to store results
	dbWriter, err := writers.NewDbWriter("sqlite://"+dbPath, false)
	if err != nil {
		return fmt.Errorf("failed to create database writer: %w", err)
	}

	// Create runner with database writer
	gwRunner, err := goWitnessRunner.NewRunner(logger, gwDriver, *options, []writers.Writer{dbWriter})
	if err != nil {
		return fmt.Errorf("failed to create gowitness runner: %w", err)
	}
	defer gwRunner.Close()

	// Feed URLs to the runner
	go func() {
		for _, url := range urls {
			gwRunner.Targets <- url
		}
		close(gwRunner.Targets)
	}()

	// Run the screenshot capture
	gwRunner.Run()

	logging.LogInfo("Gowitness completed. Screenshots saved in: " + screenshotsPath)
	logging.LogInfo("Database saved to: " + dbPath)
	logging.LogInfo("To view results, run with --server flag")
	return nil
}

func StartGoWitnessServer(opt *Options) error {
	screenshotsPath := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), "screenshots")
	dbPath := filepath.Join(screenshotsPath, "gowitness.sqlite3")

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found at %s", dbPath)
	}

	logging.LogInfo("Starting gowitness server...")
	logging.LogInfo("Database: " + dbPath)
	logging.LogInfo("Screenshots: " + screenshotsPath)
	logging.LogInfo("Server will be available at http://127.0.0.1:7171")
	logging.LogInfo("Press Ctrl+C to stop the server")

	// Import and use the gowitness web server package
	server := web.NewServer(
		"127.0.0.1",
		7171,
		"sqlite://"+dbPath,
		screenshotsPath,
	)
	server.Run()
	return nil
}
