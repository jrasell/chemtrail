package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/armon/go-metrics"
	"github.com/jrasell/chemtrail/pkg/client"
	"github.com/jrasell/chemtrail/pkg/scale"
	"github.com/jrasell/chemtrail/pkg/scale/auto"
	"github.com/jrasell/chemtrail/pkg/scale/resource"
	"github.com/jrasell/chemtrail/pkg/server/router"
	"github.com/jrasell/chemtrail/pkg/state"
	policyMemory "github.com/jrasell/chemtrail/pkg/state/policy/memory"
	scaleMemory "github.com/jrasell/chemtrail/pkg/state/scale/memory"
	"github.com/jrasell/chemtrail/pkg/watcher"
	"github.com/jrasell/chemtrail/pkg/watcher/allocs"
	"github.com/jrasell/chemtrail/pkg/watcher/nodes"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type HTTPServer struct {
	addr   string
	cfg    *Config
	logger zerolog.Logger

	// nomad is our stored Nomad client wrapper which is reused in all areas which require Nomad
	// API connectivity.
	nomad *client.Nomad

	// scaler is the backend interface used to trigger scaling.
	scaler scale.Scale

	// autoScaler is the autoscaler backend process which is responsible for triggering scaling
	// evaluations and making a decision whether scaling is required or not.
	autoscaler *auto.Scale

	// nodeWatcher is an implementation of the watcher interface used to monitor the Nomad node API
	// for changes.
	nodeWatcher watcher.Watcher

	// allocWatcher is an implementation of the watcher interface used to monitor the Nomad alloc
	// API for changes.
	allocWatcher watcher.Watcher

	// nodeResourceHandler is the handle on interacting with the Chemtrail stored node resource
	// information.
	nodeResourceHandler resource.Handler

	scaleState  state.ScaleBackend
	policyState state.PolicyBackend

	telemetry *metrics.InmemSink

	http.Server
	routes *routes
}

func New(l zerolog.Logger, cfg *Config) *HTTPServer {
	return &HTTPServer{
		addr:   fmt.Sprintf("%s:%d", cfg.Server.Bind, cfg.Server.Port),
		cfg:    cfg,
		logger: l,
		routes: &routes{},
	}
}

func (h *HTTPServer) Start() error {
	h.logger.Info().Str("addr", h.addr).Msg("starting HTTP server")

	if err := h.setup(); err != nil {
		h.logger.Error().Err(err).Msg("failed to start HTTP server")
		return err
	}

	// Start the internal node update processor.
	go h.nodeResourceHandler.RunNodeUpdateHandler()
	go h.nodeWatcher.Run(h.nodeResourceHandler.GetNodeUpdateChan())

	// Start the internal allocation update processor.
	go h.nodeResourceHandler.RunAllocUpdateHandler()
	go h.allocWatcher.Run(h.nodeResourceHandler.GetAllocUpdateChan())

	if h.cfg.Autoscale.Enabled {
		go h.autoscaler.Run()
	}

	h.handleSignals()
	return nil
}

func (h *HTTPServer) setup() error {

	if err := h.setupNomadClient(); err != nil {
		return err
	}
	h.logger.Debug().
		Str("node-id", h.nomad.NodeID).
		Msg("identified Chemtrail allocation nodeID")

	h.policyState = policyMemory.NewPolicyBackend()
	h.scaleState = scaleMemory.NewScaleStateBackend()
	h.nodeResourceHandler = resource.NewHandler(h.logger, h.nomad)

	h.scaler = scale.NewScaleBackend(&scale.BackendConfig{
		Provider:      h.cfg.Provider,
		Logger:        h.logger,
		Nomad:         h.nomad,
		NodeResources: h.nodeResourceHandler,
		ScaleState:    h.scaleState,
		PolicyState:   h.policyState},
	)

	h.nodeWatcher = nodes.NewWatcher(h.logger, h.nomad.Client)
	h.allocWatcher = allocs.NewWatcher(h.logger, h.nomad.Client)

	if h.cfg.Autoscale.Enabled {
		as, err := auto.NewAutoScaler(&auto.Config{
			Nomad:    h.nomad,
			Logger:   h.logger,
			Policy:   h.policyState,
			Resource: h.nodeResourceHandler,
			Scale:    h.scaler,
			Interval: h.cfg.Autoscale.Interval,
			Threads:  h.cfg.Autoscale.Threads,
		})
		if err != nil {
			return nil
		}
		h.autoscaler = as
	}

	// Setup telemetry based on the config passed by the operator.
	if err := h.setupTelemetry(); err != nil {
		return errors.Wrap(err, "failed to setup telemetry handler")
	}

	initialRoutes := h.setupRoutes()

	r := router.WithRoutes(h.logger, *initialRoutes)
	http.Handle("/", middlewareLogger(r, h.logger))

	// Run the TLS setup process so that if the user has configured a TLS certificate pair the
	// server uses these.
	if err := h.setupTLS(); err != nil {
		return err
	}

	// Once we have the TLS config in place, we can setup the listener which uses the TLS setup to
	// correctly start the listener.
	ln := h.setupListener()
	if ln == nil {
		return errors.New("failed to setup HTTP server, listener is nil")
	}
	h.logger.Info().Str("addr", h.addr).Msg("HTTP server successfully listening")

	go func() {
		err := h.Serve(ln)
		h.logger.Info().Err(err).Msg("HTTP server has been shutdown")
	}()

	return nil
}

func (h *HTTPServer) setupTLS() error {
	if h.cfg.TLS.CertPath != "" && h.cfg.TLS.CertKeyPath != "" {
		h.logger.Debug().Msg("setting up server TLS")

		cert, err := tls.LoadX509KeyPair(h.cfg.TLS.CertPath, h.cfg.TLS.CertKeyPath)
		if err != nil {
			return errors.Wrap(err, "failed to load certificate cert/key pair")
		}
		h.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}
	return nil
}

func (h *HTTPServer) setupNomadClient() error {
	h.logger.Debug().Msg("setting up Nomad client")

	nc, err := client.NewNomadClient()
	if err != nil {
		return err
	}
	h.nomad = nc

	return nil
}

func (h *HTTPServer) setupListener() net.Listener {
	var (
		err error
		ln  net.Listener
	)

	if h.TLSConfig != nil {
		ln, err = tls.Listen("tcp", h.addr, h.TLSConfig)
	} else {
		ln, err = net.Listen("tcp", h.addr)
	}

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to setup server HTTP listener")
	}
	return ln
}

func (h *HTTPServer) Stop() error {
	h.logger.Info().Msg("gracefully shutting down HTTP server and sub-processes")

	h.nodeResourceHandler.StopUpdateHandlers()

	// If the autoscaler is running, stop this. It is important that a Chemtrail server is given
	// time to exit cleanly as this call can take a number of seconds to complete while we
	// gracefully wait for all in-flight worker threads to finish.
	if h.autoscaler != nil && h.autoscaler.IsRunning() {
		h.autoscaler.Stop()
	}
	return h.Shutdown(context.Background())
}

func (h *HTTPServer) handleSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {

		case sig := <-signalCh:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				if err := h.Stop(); err != nil {
					h.logger.Error().Err(err).Msg("failed to cleanly shutdown server and sub-processes")
				}
				h.logger.Info().Msg("successfully shutdown server and sub-processes")
				return
			default:
				panic(fmt.Sprintf("unsupported signal: %v", sig))
			}
		}
	}
}
