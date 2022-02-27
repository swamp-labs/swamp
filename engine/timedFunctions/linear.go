package timedFunctions

import (
	"sync"
	"time"
)

type ActionParameter struct {
	TickCount int
	Tick      time.Time
	NextTick  time.Duration
}

type TickedFunction func(parameter ActionParameter) interface{}

func Constant(fct TickedFunction, rps int, tps time.Duration) chan interface{} {
	return Linear(fct, rps, rps, tps)
}

func Linear(fct TickedFunction, rpsStart, rpsEnd int, tps time.Duration) chan interface{} {
	// Calculate the coefficient
	startInterval := time.Second / time.Duration(rpsStart)
	endInterval := time.Second / time.Duration(rpsEnd)
	deltaInterval := float64(endInterval.Milliseconds() - startInterval.Milliseconds())
	linearCoefficient := deltaInterval / float64(tps.Milliseconds())

	// Results channel will contain the results of all requests
	results := make(chan interface{})
	// Done channel will stop the main loop that send requests
	done := make(chan bool)

	ticker := time.NewTicker(startInterval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	tickIndex := 0
	go func() {
		var requestsWaitingGroup sync.WaitGroup
		defer func() {
			requestsWaitingGroup.Wait()
			wg.Done()
		}()

		startTime := time.Now()
		for {
			// select allows to listen to channels
			select {
			case <-done:
				// Stop the loop
				return
			case t := <-ticker.C:
				// New tick ! Run function
				requestsWaitingGroup.Add(1)
				nextTickDuration := time.Duration(int64(float64(t.Sub(startTime).Milliseconds())*linearCoefficient))*time.Millisecond + startInterval
				go func() {
					defer requestsWaitingGroup.Done()
					results <- fct(ActionParameter{
						TickCount: tickIndex,
						Tick:      t,
						NextTick:  nextTickDuration,
					})
				}()
				// Set the new ticker time
				ticker.Reset(nextTickDuration)
			}
			tickIndex += 1
		}
	}()

	// Wait until the end of the timer
	time.Sleep(tps)

	// Stop the ticker
	ticker.Stop()
	done <- true

	// Wait for the last request before closing the results
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
