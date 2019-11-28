package allocs

import (
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/chemtrail/pkg/watcher"
	"github.com/rs/zerolog"
)

type Watcher struct {
	logger          zerolog.Logger
	nomad           *api.Client
	lastChangeIndex uint64
}

func NewWatcher(logger zerolog.Logger, nomad *api.Client) watcher.Watcher {
	return &Watcher{
		logger: logger,
		nomad:  nomad,
	}
}

func (w *Watcher) Run(updateChan chan interface{}) {
	w.logger.Info().Msg("starting Chemtrail Nomad alloc watcher")

	var maxFound uint64

	q := &api.QueryOptions{WaitTime: 5 * time.Minute, WaitIndex: 1}

	for {

		allocs, meta, err := w.nomad.Allocations().List(q)
		if err != nil {
			w.logger.Error().Err(err).Msg("failed to call Nomad API for alloc listing")
			time.Sleep(10 * time.Second)
			continue
		}

		if !watcher.IndexHasChange(meta.LastIndex, q.WaitIndex) {
			w.logger.Debug().Msg("alloc watcher last index has not changed")
			continue
		}
		w.logger.Debug().
			Uint64("old", q.WaitIndex).
			Uint64("new", meta.LastIndex).
			Msg("alloc watcher last index has changed")

		// Iterate over all the returned node allocs.
		for i := range allocs {

			if !watcher.IndexHasChange(allocs[i].ModifyIndex, w.lastChangeIndex) {
				continue
			}

			w.logger.Debug().
				Uint64("old", w.lastChangeIndex).
				Uint64("new", allocs[i].ModifyIndex).
				Msg("alloc modify index has changed is greater than last recorded")

			alloc, _, err := w.nomad.Allocations().Info(allocs[i].ID, nil)
			if err != nil {
				w.logger.Error().
					Str("alloc-id", allocs[i].ID).
					Err(err).
					Msg("failed to call Nomad API for alloc info")
			}

			maxFound = watcher.MaxFound(allocs[i].ModifyIndex, maxFound)
			if alloc != nil {
				updateChan <- alloc
			}
		}

		// Update the Nomad API wait index to start long polling from the correct point and update
		// our recorded lastChangeIndex so we have the correct point to use during the next API
		// return.
		q.WaitIndex = meta.LastIndex
		w.lastChangeIndex = maxFound
	}
}
