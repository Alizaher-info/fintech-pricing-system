package worker

import (
	"fmt"
	"time"

	"services/pricing-fetcher/internal/fetcher"
	"services/pricing-fetcher/pkg/logger"
)

type PriceWorker struct {
	priceFetcher *fetcher.PriceFetcher
	interval     int
	stopChan     chan bool
}

func NewPriceWorker(priceFetcher *fetcher.PriceFetcher, intervalSeconds int) *PriceWorker {
	return &PriceWorker{
		priceFetcher: priceFetcher,
		interval:     intervalSeconds,
		stopChan:     make(chan bool),
	}
}

// Start begins the background worker that fetches prices periodically
func (pw *PriceWorker) Start() {
	go func() {
		// Initial fetch on startup
		logger.Info("Performing initial price fetch...")
		if err := pw.priceFetcher.FetchAndSaveAll(); err != nil {
			logger.Error(fmt.Sprintf("Initial fetch failed: %v", err))
		}

		// Create ticker for periodic fetching
		ticker := time.NewTicker(time.Duration(pw.interval) * time.Second)
		defer ticker.Stop()

		logger.Info(fmt.Sprintf("Background worker started - fetching every %d seconds", pw.interval))

		for {
			select {
			case <-ticker.C:
				// Fetch prices on each tick
				logger.Info("Ticker triggered - fetching prices...")
				if err := pw.priceFetcher.FetchAndSaveAll(); err != nil {
					logger.Error(fmt.Sprintf("Fetch failed: %v", err))
				}

			case <-pw.stopChan:
				// Stop signal received
				logger.Info("Background worker stopping...")
				return
			}
		}
	}()
}

// Stop gracefully stops the background worker
func (pw *PriceWorker) Stop() {
	logger.Info("Stopping background worker...")
	close(pw.stopChan)
	
	// Give goroutine time to finish gracefully
	time.Sleep(2 * time.Second)
	
	logger.Info("Background worker stopped.")
}
