package runner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/henrique-gomesz/joeyscan4me/pkg/logging"
)

var SubfinderOutputFile = "subdomains.txt"
var HttpxOutputFile = "up_subdomains.txt"
var HttpxTechOutputFile = "up_subdomains_with_tech.txt"
var KatanaOutputFile = "crawling_results.txt"

func GetOutputFilePath(workdir, domain string) string {
	return filepath.Join(workdir, "output", domain)
}

func CreateOutputFile(filePath string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		logging.LogError("Failed to create output directory", err)
		return nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		logging.LogError("Failed to create output file", err)
		return nil, err
	}

	return file, nil
}

func OpenOutputFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logging.LogError("Failed to open file", err)
		return nil, err
	}
	return file, nil
}

func ReadFileLines(filePath string) ([]string, error) {
	file, err := OpenOutputFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logging.LogError("Failed to read file", err)
		return nil, err
	}

	return lines, nil
}

func NormalizeAndDedupeLines(lines []string) []string {
	seen := make(map[string]struct{}, len(lines))
	normalized := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if _, exists := seen[line]; exists {
			continue
		}

		seen[line] = struct{}{}
		normalized = append(normalized, line)
	}

	return normalized
}

func WriteToFile(filePath string, content string) error {
	file, err := CreateOutputFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		logging.LogError("Failed to write to file", err)
		return err
	}

	return nil
}

func AppendToFile(filePath string, content string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		logging.LogError("Failed to create output directory", err)
		return err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.LogError("Failed to open file for appending", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		logging.LogError("Failed to append to file", err)
		return err
	}

	return nil
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FileNonEmpty returns true when filePath exists, is a regular file, and has size > 0.
func FileNonEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir() && info.Size() > 0
}
