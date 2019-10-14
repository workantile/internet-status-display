package main

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/ddo/go-fast.v0"
)

func startFastcomSpeedTest(interval time.Duration, cb func(kbps float64, err error)) {
	scheduleAfter(func() {
		fastCom := fast.New()
		err := fastCom.Init()
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed to initialize"))
			return
		}

		urls, err := fastCom.GetUrls()
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed to get URLs"))
			return
		}

		kbpsChan := make(chan float64)
		go func() {
			for kbps := range kbpsChan {
				cb(kbps, nil)
			}
		}()

		ProbeMutex.Lock()
		defer ProbeMutex.Unlock()
		err = fastCom.Measure(urls, kbpsChan)
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed"))
			return
		}
	}, interval)
}