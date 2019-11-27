package nodes

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
	w.logger.Info().Msg("starting Chemtrail Nomad nodes watcher")

	var maxFound uint64

	q := &api.QueryOptions{WaitTime: 5 * time.Minute, WaitIndex: 1}

	for {

		nodes, meta, err := w.nomad.Nodes().List(q)
		if err != nil {
			w.logger.Error().Err(err).Msg("failed to call Nomad API for nodes listing")
			time.Sleep(10 * time.Second)
			continue
		}

		if !watcher.IndexHasChange(meta.LastIndex, q.WaitIndex) {
			w.logger.Debug().Msg("nodes watcher last index has not changed")
			continue
		}
		w.logger.Debug().
			Uint64("old", q.WaitIndex).
			Uint64("new", meta.LastIndex).
			Msg("nodes watcher last index has changed")

		// Iterate over all the returned nodes.
		for i := range nodes {

			if !watcher.IndexHasChange(nodes[i].ModifyIndex, w.lastChangeIndex) {
				continue
			}

			w.logger.Debug().
				Uint64("old", w.lastChangeIndex).
				Uint64("new", nodes[i].ModifyIndex).
				Str("node", nodes[i].ID).
				Msg("node modify index has changed is greater than last recorded")

			node, _, err := w.nomad.Nodes().Info(nodes[i].ID, nil)
			if err != nil {
				w.logger.Error().
					Str("node-id", nodes[i].ID).
					Err(err).
					Msg("failed to call Nomad API for node info")
			}

			maxFound = watcher.MaxFound(nodes[i].ModifyIndex, maxFound)
			updateChan <- node
		}

		// Update the Nomad API wait index to start long polling from the correct point and update
		// our recorded lastChangeIndex so we have the correct point to use during the next API
		// return.
		q.WaitIndex = meta.LastIndex
		w.lastChangeIndex = maxFound
	}
}
