package main

import (
	"fmt"
	"time"
)

func scheduleAfter(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("[critical] panic in scheduled goroutine", r)
			}
		}()
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()
	return stop
}
