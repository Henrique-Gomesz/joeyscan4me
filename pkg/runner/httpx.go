package runner

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/henrique-gomesz/joeyscan4me/pkg/logging"

	httpRunner "github.com/projectdiscovery/httpx/runner"
)

func RunHttpx(opt *Options) error {
	outputPath := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), HttpxOutputFile)
	techOutputPath := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), HttpxTechOutputFile)
	subfinderFile := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), SubfinderOutputFile)

	if opt.Resume && FileNonEmpty(outputPath) {
		logging.LogInfo("Skipping httpx — output already exists: " + outputPath)
		return nil
	}

	targets, err := ReadFileLines(subfinderFile)
	if err != nil {
		return fmt.Errorf("failed to read subfinder output file: %w", err)
	}
	targets = NormalizeAndDedupeLines(targets)

	file, err := CreateOutputFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create httpx output file: %w", err)
	}
	defer file.Close()

	techFile, err := CreateOutputFile(techOutputPath)
	if err != nil {
		return fmt.Errorf("failed to create httpx tech output file: %w", err)
	}
	defer techFile.Close()

	if len(targets) == 0 {
		logging.LogInfo("No subdomains found to probe with httpx")
		return nil
	}

	var writeMu sync.Mutex
	seen := make(map[string]struct{}, len(targets))

	httpxOpts := &httpRunner.Options{
		Methods:            "GET",
		FollowRedirects:    true,
		TechDetect:         true,
		RandomAgent:        true,
		InputFile:          subfinderFile,
		Threads:            opt.HttpxThreads,
		RateLimit:          opt.HttpxRateLimit,
		Timeout:            opt.HttpxTimeout,
		StatusCode:         true,
		ExtractTitle:       true,
		Location:           true,
		OutputServerHeader: true,
		OnResult: func(r httpRunner.Result) {
			if r.Err != nil {
				logging.LogError("HTTPX error", r.Err)
				return
			}

			resultURL := strings.TrimSpace(r.URL)
			if resultURL == "" {
				return
			}

			writeMu.Lock()
			defer writeMu.Unlock()

			if _, exists := seen[resultURL]; exists {
				return
			}
			seen[resultURL] = struct{}{}

			if _, err := fmt.Fprintf(file, "%s\n", resultURL); err != nil {
				logging.LogError("Failed to write URL to httpx output", err)
				return
			}

			if len(r.Technologies) > 0 {
				techList := strings.Join(r.Technologies, ", ")
				if _, err := fmt.Fprintf(techFile, "%s [%s]\n", resultURL, techList); err != nil {
					logging.LogError("Failed to write technologies to httpx tech output", err)
				}
				return
			}

			if _, err := fmt.Fprintf(techFile, "%s [no technologies detected]\n", resultURL); err != nil {
				logging.LogError("Failed to write default technology marker", err)
			}
		},
	}

	if strings.TrimSpace(opt.HttpxPorts) != "" {
		if err := httpxOpts.CustomPorts.Set(opt.HttpxPorts); err != nil {
			return fmt.Errorf("invalid --httpx-ports value: %w", err)
		}
	}

	if err := httpxOpts.ValidateOptions(); err != nil {
		return fmt.Errorf("failed to validate httpx options: %w", err)
	}

	httpx, err := httpRunner.New(httpxOpts)
	if err != nil {
		return fmt.Errorf("failed to create httpx runner: %w", err)
	}

	defer httpx.Close()

	logging.LogInfo(fmt.Sprintf("Running httpx (threads=%d rate=%d timeout=%ds)", opt.HttpxThreads, opt.HttpxRateLimit, opt.HttpxTimeout))
	httpx.RunEnumeration()

	if len(seen) == 0 {
		logging.LogInfo("HTTPX finished with no live services")
		return nil
	}

	logging.LogSuccess(fmt.Sprintf("HTTPX finished with %d unique live URLs", len(seen)))
	return nil
}
