package runner

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/henrique-gomesz/joeyscan4me/pkg/logging"
)

func StartScan(opt *Options) error {
	startedAt := time.Now()

	if opt.Resume {
		logging.LogInfo("Resume mode enabled — stages with existing output will be skipped")
	}

	if err := RunSubfinder(opt); err != nil {
		return fmt.Errorf("subfinder failed: %w", err)
	}

	if err := RunHttpx(opt); err != nil {
		return fmt.Errorf("httpx failed: %w", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// in order to speed up things we run katana and gowitness concurrently
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := RunKatana(opt); err != nil {
			errCh <- fmt.Errorf("katana failed: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := RunGowitness(opt); err != nil {
			errCh <- fmt.Errorf("gowitness failed: %w", err)
		}
	}()

	wg.Wait()
	close(errCh)

	var postProcessingErrors []error
	for err := range errCh {
		postProcessingErrors = append(postProcessingErrors, err)
	}

	if len(postProcessingErrors) > 0 {
		return errors.Join(postProcessingErrors...)
	}

	if opt.SummaryJSON {
		summaryPath, err := WriteScanSummary(opt, startedAt, time.Now())
		if err != nil {
			return fmt.Errorf("failed to write scan summary: %w", err)
		}
		logging.LogInfo("Scan summary written to " + summaryPath)
	}

	// start go witness server if --server flag is set
	if opt.Server {
		if err := StartGoWitnessServer(opt); err != nil {
			return fmt.Errorf("gowitness server failed: %w", err)
		}
	}

	return nil
}
