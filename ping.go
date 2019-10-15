package main

import (
	"log"
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

const pingCount = 5
const pingTimeout = 3 * time.Second

func startPing(interval time.Duration, targetAddr string, cb func(stats *PingResult, err error)) {
	log.Printf("[info] scheduling ping to %s, every %d seconds...\n", targetAddr, int(interval.Seconds()))
	scheduleAfter(func() {
		pinger, err := ping.NewPinger(targetAddr)
		if err != nil {
			cb(nil, errors.Wrap(err, "building new Pinger failed"))
			return
		}
		pinger.Count = pingCount
		pinger.Timeout = pingTimeout

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
