package worker

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"code.cloudfoundry.org/lager/v3"
	"code.cloudfoundry.org/lager/v3/lagerctx"
	"github.com/concourse/concourse/tsa"
)

type Beacon struct {
	Logger lager.Logger

	Client TSAClient

	DrainSignals <-chan os.Signal

	RebalanceInterval      time.Duration
	ConnectionDrainTimeout time.Duration

	LocalGardenNetwork string
	LocalGardenAddr    string

	LocalBaggageclaimNetwork string
	LocalBaggageclaimAddr    string

	drained int32
}

// total number of active registrations; all but one are "live", the rest
// should all be draining
const maxActiveRegistrations = 5

func (beacon *Beacon) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	beacon.Logger.Debug("start")
	defer beacon.Logger.Debug("done")

	var rebalanceCh <-chan time.Time
	if beacon.RebalanceInterval != 0 {
		ticker := time.NewTicker(beacon.RebalanceInterval)
		defer ticker.Stop()

		rebalanceCh = ticker.C
	}

	cwg := &countingWaitGroup{}
	defer cwg.Wait()

	rootCtx, cancelAll := context.WithCancel(lagerctx.NewContext(context.Background(), beacon.Logger))
	defer cancelAll()

	latestErrChan := make(chan error, 1)

	// Make sure we are concurrently handling signals while registering a worker since worker registration can block
	// indefinitely.
	receivedSignals := make(chan os.Signal, 1)
	go func() {
		select {
		case s := <-signals:
			// Cancel context to unblock worker registration
			cancelAll()
			// Pass on signal to the main loop that runs after worker registration
			receivedSignals <- s
		case <-rootCtx.Done():
		}
	}()

	cwg.Add(1)

	ctx, cancel := context.WithCancel(rootCtx)
	beacon.registerWorker(ctx, cwg, latestErrChan)

	close(ready)

	var retiring bool

	cancelPrev := cancel
	for {
		select {
		case <-rebalanceCh:
			logger := beacon.Logger.Session("rebalance")

			if cwg.Count() >= maxActiveRegistrations {
				logger.Info("max-active-registrations-reached", lager.Data{
					"limit": maxActiveRegistrations,
				})

				continue
			} else {
				logger.Debug("rebalancing")
			}

			ctx, cancel := context.WithCancel(lagerctx.NewContext(rootCtx, logger))

			// make a new channel so prior registrations can write to their own
			// buffered channel and exit
			latestErrChan = make(chan error, 1)

			cwg.Add(1)
			beacon.registerWorker(ctx, cwg, latestErrChan)

			cancelPrev()
			cancelPrev = cancel

		case err := <-latestErrChan:
			if err != nil {
				beacon.Logger.Error("exited-with-error", err)
			} else {
				beacon.Logger.Info("exited")
			}

			return err

		case sig := <-beacon.DrainSignals:
			atomic.StoreInt32(&beacon.drained, 1)

			logger := beacon.Logger.Session("drain")

			logger.Debug("received-drain-signal", lager.Data{
				"signal": sig.String(),
			})

			// prevent rebalancing from switching the worker back to 'running'
			rebalanceCh = nil

			if isLand(sig) {
				logger.Info("landing-worker")

				err := beacon.Client.Land(rootCtx)
				if err != nil {
					logger.Error("failed-to-land-worker", err)

					return err
				}
			} else if isRetire(sig) {
				retiring = true

				logger.Info("retiring-worker")

				err := beacon.Client.Retire(rootCtx)
				if err != nil {
					logger.Error("failed-to-retire-worker", err)

					return err
				}
			}

		case <-receivedSignals:
			logger := beacon.Logger.Session("signal")

			logger.Info("signalled")

			if retiring {
				logger.Info("deleting-worker")

				err := beacon.Client.Delete(rootCtx)
				if err != nil {
					logger.Error("failed-to-delete-worker", err)
					return err
				}
			}

			return nil
		}
	}
}

func (beacon *Beacon) Drained() bool {
	return atomic.LoadInt32(&beacon.drained) == 1
}

func (beacon *Beacon) registerWorker(
	ctx context.Context,
	cwg *countingWaitGroup,
	errs chan<- error,
) {
	logger := lagerctx.FromContext(ctx)

	once := &sync.Once{}

	registeredOrFailed := make(chan struct{})
	go func() {
		defer cwg.Done()

		errs <- beacon.Client.Register(ctx, tsa.RegisterOptions{
			LocalGardenNetwork: beacon.LocalGardenNetwork,
			LocalGardenAddr:    beacon.LocalGardenAddr,

			LocalBaggageclaimNetwork: beacon.LocalBaggageclaimNetwork,
			LocalBaggageclaimAddr:    beacon.LocalBaggageclaimAddr,

			ConnectionDrainTimeout: beacon.ConnectionDrainTimeout,

			RegisteredFunc: func() {
				logger.Info("registered")
				once.Do(func() { close(registeredOrFailed) })
			},

			HeartbeatedFunc: func() {
				logger.Debug("heartbeated")
			},
		})

		once.Do(func() { close(registeredOrFailed) })
	}()

	<-registeredOrFailed
}
