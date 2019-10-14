package main

import (
	"math"
	"time"

	"github.com/pkg/errors"
	"github.com/sparrc/go-ping"
)

type PingResult struct {
	AvgRtt        int64 `json:"avgRtt"`
	MaxRtt        int64 `json:"maxRtt"`
	StdDevRtt     int64 `json:"stdDevRtt"`
	PacketLossPct int   `json:"packetLossPct"`
}

func startPing(interval time.Duration, targetAddr string, cb func(stats *PingResult, err error)) {
	scheduleAfter(func() {
		pinger, err := ping.NewPinger(targetAddr)
		if err != nil {
			cb(nil, errors.Wrap(err, "building new Pinger failed"))
			return
		}
		pinger.Count = 5
		pinger.Timeout = 2 * time.Second

		ProbeMutex.Lock()
		pinger.Run()
		ProbeMutex.Unlock()

		pingStats := pinger.Statistics()
		cb(&PingResult{
			AvgRtt:        pingStats.AvgRtt.Milliseconds(),
			MaxRtt:        pingStats.MaxRtt.Milliseconds(),
			StdDevRtt:     pingStats.StdDevRtt.Milliseconds(),
			PacketLossPct: int(math.Round(pingStats.PacketLoss)),
		}, nil)
	}, interval)
}
