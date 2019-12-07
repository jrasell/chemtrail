package server

import "time"

// gcEvalPeriod is the time period at which the automatic scaling state garbage collector is run
// at.
var gcEvalPeriod = time.Minute * 10

// runGarbageCollectionLoop is responsible for periodically running the scaling state garbage
// collection function.
func (h *HTTPServer) runGarbageCollectionLoop() {
	h.logger.Info().Msg("started scaling state garbage collector handler")

	h.gcIsRunning = true

	t := time.NewTicker(gcEvalPeriod)
	defer t.Stop()

	for {
		select {
		case <-h.stopChan:
			h.logger.Info().Msg("shutting down state garbage collection handler")
			h.gcIsRunning = false
			return
		case <-t.C:
			h.logger.Debug().Msg("triggering internal run of state garbage collection")
			h.scaleState.RunStateGarbageCollection()
		}
	}
}
