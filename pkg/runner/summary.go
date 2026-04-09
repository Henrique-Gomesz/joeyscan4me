package runner

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type scanSummary struct {
	Domain          string        `json:"domain"`
	Profile         string        `json:"profile"`
	Workdir         string        `json:"workdir"`
	OutputDirectory string        `json:"output_directory"`
	StartedAt       string        `json:"started_at"`
	FinishedAt      string        `json:"finished_at"`
	DurationSeconds float64       `json:"duration_seconds"`
	Counts          summaryCounts `json:"counts"`
	Files           summaryFiles  `json:"files"`
}

type summaryCounts struct {
	Subdomains       int `json:"subdomains"`
	LiveURLs         int `json:"live_urls"`
	LiveURLsWithTech int `json:"live_urls_with_tech"`
	CrawledEntries   int `json:"crawled_entries"`
	Screenshots      int `json:"screenshots"`
}

type summaryFiles struct {
	SubdomainsFile        string `json:"subdomains_file"`
	LiveURLsFile          string `json:"live_urls_file"`
	LiveURLsWithTechFile  string `json:"live_urls_with_tech_file"`
	KatanaOutputFile      string `json:"katana_output_file"`
	GowitnessDatabaseFile string `json:"gowitness_database_file"`
}

func WriteScanSummary(opt *Options, startedAt, finishedAt time.Time) (string, error) {
	outputDir := GetOutputFilePath(opt.Workdir, opt.Domain)
	summaryPath := filepath.Join(outputDir, opt.SummaryFile)

	subdomainsFile := filepath.Join(outputDir, SubfinderOutputFile)
	liveURLsFile := filepath.Join(outputDir, HttpxOutputFile)
	liveURLsWithTechFile := filepath.Join(outputDir, HttpxTechOutputFile)
	katanaOutputFile := filepath.Join(outputDir, KatanaOutputFile)
	screenshotsDir := filepath.Join(outputDir, "screenshots")
	gowitnessDBFile := filepath.Join(screenshotsDir, "gowitness.sqlite3")

	subdomainsCount, err := countUniqueLinesFromFile(subdomainsFile)
	if err != nil {
		return "", err
	}
	liveURLsCount, err := countUniqueLinesFromFile(liveURLsFile)
	if err != nil {
		return "", err
	}
	liveURLsWithTechCount, err := countUniqueLinesFromFile(liveURLsWithTechFile)
	if err != nil {
		return "", err
	}
	katanaEntriesCount, err := countUniqueLinesFromFile(katanaOutputFile)
	if err != nil {
		return "", err
	}
	screenshotCount, err := countScreenshotFiles(screenshotsDir)
	if err != nil {
		return "", err
	}

	summary := scanSummary{
		Domain:          opt.Domain,
		Profile:         opt.Profile,
		Workdir:         opt.Workdir,
		OutputDirectory: outputDir,
		StartedAt:       startedAt.Format(time.RFC3339),
		FinishedAt:      finishedAt.Format(time.RFC3339),
		DurationSeconds: finishedAt.Sub(startedAt).Seconds(),
		Counts: summaryCounts{
			Subdomains:       subdomainsCount,
			LiveURLs:         liveURLsCount,
			LiveURLsWithTech: liveURLsWithTechCount,
			CrawledEntries:   katanaEntriesCount,
			Screenshots:      screenshotCount,
		},
		Files: summaryFiles{
			SubdomainsFile:        subdomainsFile,
			LiveURLsFile:          liveURLsFile,
			LiveURLsWithTechFile:  liveURLsWithTechFile,
			KatanaOutputFile:      katanaOutputFile,
			GowitnessDatabaseFile: gowitnessDBFile,
		},
	}

	content, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(summaryPath), 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(summaryPath, append(content, '\n'), 0644); err != nil {
		return "", err
	}

	return summaryPath, nil
}

func countUniqueLinesFromFile(filePath string) (int, error) {
	if !FileExists(filePath) {
		return 0, nil
	}

	lines, err := ReadFileLines(filePath)
	if err != nil {
		return 0, err
	}

	return len(NormalizeAndDedupeLines(lines)), nil
}

func countScreenshotFiles(root string) (int, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return 0, nil
	}

	count := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
