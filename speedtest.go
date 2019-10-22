package main

import (
	"log"
	"time"

	"github.com/kylegrantlucas/speedtest"
	"github.com/pkg/errors"
	"gopkg.in/ddo/go-fast.v0"
)

func startFastComSpeedTest(interval time.Duration, cb func(kbps float64, err error)) {
	log.Printf("[info] scheduling fast.com speedtest, every %d seconds...\n", int(interval.Seconds()))
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

func startSpeedtestNetUploadSpeedTest(interval time.Duration, cb func(kbps float64, err error)) {
	log.Printf("[info] scheduling speedtest.net upload test, every %d seconds...\n", int(interval.Seconds()))
	scheduleAfter(func() {
		client, err := speedtest.NewDefaultClient()
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed to initialize"))
			return
		}

		server, err := client.GetServer("")
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed to get server"))
			return
		}

		ProbeMutex.Lock()
		upMbps, err := client.Upload(server)
		ProbeMutex.Unlock()
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed"))
			return
		}

		upKbps := upMbps * 1000
		cb(upKbps, nil)
	}, interval)
}

func startSpeedtestNetDownloadSpeedTest(interval time.Duration, cb func(kbps float64, err error)) {
	log.Printf("[info] scheduling speedtest.net download test, every %d seconds...\n", int(interval.Seconds()))
	scheduleAfter(func() {
		client, err := speedtest.NewDefaultClient()
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed to initialize"))
			return
		}

		server, err := client.GetServer("")
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed to get server"))
			return
		}

		ProbeMutex.Lock()
		downMbps, err := client.Download(server)
		ProbeMutex.Unlock()
		if err != nil {
			cb(0, errors.Wrap(err, "scheduled speedtest failed"))
			return
		}

		downKbps := downMbps * 1000
		cb(downKbps, nil)
	}, interval)
}
