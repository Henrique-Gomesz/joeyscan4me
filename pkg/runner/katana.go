package runner

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/Henrique-Gomesz/JoeyScan4Me/pkg/logging"

	"github.com/projectdiscovery/katana/pkg/engine/standard"
	katanaTypes "github.com/projectdiscovery/katana/pkg/types"
)

func RunKatana(opt *Options) error {
	httpxOutputPath := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), HttpxOutputFile)
	katanaOutputPath := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), KatanaOutputFile)

	urls, err := ReadFileLines(httpxOutputPath)
	if err != nil {
		return fmt.Errorf("failed to read httpx output file: %w", err)
	}
	urls = NormalizeAndDedupeLines(urls)

	if len(urls) == 0 {
		logging.LogInfo("No URLs found to crawl")
		emptyFile, err := CreateOutputFile(katanaOutputPath)
		if err != nil {
			return fmt.Errorf("failed to create empty katana output file: %w", err)
		}
		emptyFile.Close()
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(katanaOutputPath), 0755); err != nil {
		return fmt.Errorf("failed to create katana output directory: %w", err)
	}

	logging.LogInfo("Starting crawling with Katana")
	logging.LogInfo(fmt.Sprintf("Crawling %d URLs", len(urls)))

	katanaOpts := &katanaTypes.Options{
		URLs:                   urls,
		MaxDepth:               opt.KatanaDepth,
		BodyReadSize:           math.MaxInt,
		Timeout:                opt.KatanaTimeout,
		Concurrency:            opt.KatanaConcurrency,
		Parallelism:            opt.KatanaParallelism,
		Retries:                3,
		RateLimit:              opt.KatanaRateLimit,
		Strategy:               "depth-first",
		ScrapeJSResponses:      true,
		ScrapeJSLuiceResponses: true,
		OutputFile:             katanaOutputPath,
		ExtensionFilter:        []string{"css"},
	}

	crawlerOptions, err := katanaTypes.NewCrawlerOptions(katanaOpts)
	if err != nil {
		return fmt.Errorf("failed to create katana crawler options: %w", err)
	}
	defer crawlerOptions.Close()

	crawler, err := standard.New(crawlerOptions)
	if err != nil {
		return fmt.Errorf("failed to create katana crawler: %w", err)
	}
	defer crawler.Close()

	failed := 0
	for _, targetURL := range urls {
		if err := crawler.Crawl(targetURL); err != nil {
			failed++
			logging.LogError(fmt.Sprintf("Could not crawl %s", targetURL), err)
		}
	}

	if failed > 0 {
		return fmt.Errorf("katana failed to crawl %d/%d URLs", failed, len(urls))
	}

	logging.LogSuccess("Results saved to " + katanaOutputPath)
	return nil
}
