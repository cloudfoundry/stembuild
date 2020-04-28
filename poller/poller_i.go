package poller

import "time"

//go:generate counterfeiter . PollerI
type PollerI interface {
	Poll(duration time.Duration, loopFunc func() (bool, error)) error
}
