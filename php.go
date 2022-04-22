package caddyphpfpm

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

var emptyFunc = func() {}

const (
	minRestartDelay             = time.Duration(0)
	maxRestartDelay             = 5 * time.Minute
	terminationGracePeriod      = 15 * time.Second
	durationToResetRestartDelay = 10 * time.Minute
)

// PHP provides functionality to start and supervise php-fpm in the background
type PHP struct {
	ctx       context.Context
	cancelCtx context.CancelFunc

	cmd         *exec.Cmd
	waiter      sync.WaitGroup
	keepRunning bool

	command      []string
	sockLocation string
	startTimeout time.Duration
}

func NewPHP() *PHP {
	ctx, cancel := context.WithCancel(context.Background())
	return &PHP{ctx: ctx, cancelCtx: cancel, command: []string{"php-fpm"}, sockLocation: "fpm.sock", startTimeout: 10 * time.Second}
}

func (p *PHP) wait() error {
	ctx, cancel := context.WithDeadline(p.ctx, time.Now().Add(p.startTimeout))
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return errors.New("Failed to wait for php to start")
		case <-time.After(1 * time.Second):
			if _, err := os.Stat(p.sockLocation); err == nil {
				L.Info("Found php socket", zap.String("sock_location", p.sockLocation))
				return nil
			}
			L.Warn("Still haven't seen php socket", zap.String("sock_location", p.sockLocation))
		}
	}
}

// Run a process and supervise
func (p *PHP) Run() {
	p.keepRunning = true

	restartDelay := minRestartDelay

	for p.keepRunning {
		L.Info("Starting php-fpm", zap.String("command", strings.Join(p.command, " ")))
		p.cmd = exec.Command(p.command[0], p.command[1:]...)
		p.cmd.Stdout = os.Stderr
		p.cmd.Stderr = os.Stderr
		configureSysProcAttr(p.cmd)
		configureExecutingUser(p.cmd, "")

		start := time.Now()
		err := p.cmd.Start()

		failedAtStart := false
		if err != nil {
			L.Error("failed to start process", zap.Error(err))
			failedAtStart = true
		} else {
			L.Info("process started", zap.Int("pid", p.cmd.Process.Pid))

			p.waiter.Add(1)
			err = p.cmd.Wait()
			p.waiter.Done()
		}

		duration := time.Now().Sub(start)

		if err != nil {
			if !failedAtStart {
				L.Error("process exited with error", zap.Error(err), zap.Duration("duration", duration))
			}
		} else {
			L.Info("process exited", zap.Duration("duration", duration))
		}

		if !p.keepRunning {
			break
		}

		if p.keepRunning {
			if restartDelay > minRestartDelay && (err == nil || duration > durationToResetRestartDelay) {
				L.Info("resetting restart delay", zap.Duration("delay", minRestartDelay))
				restartDelay = minRestartDelay
			}

			if err != nil {
				L.Info("process will restart", zap.Duration("wait_delay", restartDelay))
				time.Sleep(restartDelay)
				restartDelay = increaseRestartDelay(restartDelay)
			}
		}
	}
}

// Stop the supervised process
func (p *PHP) Stop() {
	p.cancelCtx()
	p.keepRunning = false

	if cmdIsRunning(p.cmd) {
		L.Debug("sending 'interrupt signal to gracefully stop the process")

		err := p.cmd.Process.Signal(os.Interrupt)
		if err == nil {
			go func() {
				time.Sleep(terminationGracePeriod)
				if cmdIsRunning(p.cmd) {
					L.Info("termination grace period exceeded, killing")

					p.cmd.Process.Kill()
				}
			}()

			p.waiter.Wait()
		} else {
			L.
				With(zap.Error(err)).
				Info("error while sending 'interupt' signal, killing")

			p.cmd.Process.Kill()
		}
	}
}

func cmdIsRunning(cmd *exec.Cmd) bool {
	return cmd != nil && cmd.Process != nil && (cmd.ProcessState == nil || !cmd.ProcessState.Exited())
}

func increaseRestartDelay(restartDelay time.Duration) time.Duration {
	if restartDelay == 0 {
		return 1 * time.Second
	}

	restartDelay = restartDelay * 2

	if restartDelay > maxRestartDelay {
		restartDelay = maxRestartDelay
	}

	return restartDelay
}
