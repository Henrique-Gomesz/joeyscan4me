package runner

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Henrique-Gomesz/JoeyScan4Me/pkg/logging"

	subfinderRunner "github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

func RunSubfinder(opt *Options) error {
	filePath := filepath.Join(GetOutputFilePath(opt.Workdir, opt.Domain), SubfinderOutputFile)
	file, err := CreateOutputFile(filePath)

	if err != nil {
		return fmt.Errorf("failed to create subfinder output file: %w", err)
	}

	defer file.Close()

	subfinderOpts := &subfinderRunner.Options{
		Threads:            opt.SubfinderThreads,
		Timeout:            opt.SubfinderTimeout,
		MaxEnumerationTime: opt.SubfinderMaxTime,
		All:                true,
		Domain:             []string{opt.Domain},
		OutputFile:         filePath,
		OutputDirectory:    filepath.Dir(filePath),
		Output:             file,
	}

	subfinder, err := subfinderRunner.NewRunner(subfinderOpts)
	if err != nil {
		return fmt.Errorf("failed to create subfinder runner: %w", err)
	}

	logging.LogInfo("Running subfinder")
	if err = subfinder.RunEnumerationWithCtx(context.Background()); err != nil {
		return fmt.Errorf("failed to enumerate subdomains: %w", err)
	}

	logging.LogSuccess("Subfinder finished")
	return nil
}
