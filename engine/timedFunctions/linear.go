package timedFunctions

import (
	"sync"
	"time"
)

type TickFunctionParameter struct {
	TickCount int
	Tick      time.Time
	NextTick  time.Duration
}

type TickedFunction func(parameter TickFunctionParameter) interface{}
type tickVariationFunction func(elapsedTime time.Duration) time.Duration

// Constant calls a function multiple times during an interval with a frequency of rps.
func Constant(fct TickedFunction, rps int, tps time.Duration) chan interface{} {
	return Linear(fct, rps, rps, tps)
}

// makeLinearVariationFunction creates a variation function from the starting and ending rps during a given time
func makeLinearVariationFunction(rpsStart, rpsEnd int, tps time.Duration) tickVariationFunction {
	startInterval := time.Second / time.Duration(rpsStart)
	endInterval := time.Second / time.Duration(rpsEnd)
	deltaInterval := float64(endInterval.Milliseconds() - startInterval.Milliseconds())
	linearCoefficient := deltaInterval / float64(tps.Milliseconds())

	return func(elapsedTime time.Duration) time.Duration {
		return time.Duration(float64(elapsedTime.Milliseconds())*linearCoefficient)*time.Millisecond + startInterval
	}
}

// tickFunction initialize a tick loop and returns a channel that contains all results for each function calls.
func tickFunction(tickedFunction TickedFunction, timeFunction tickVariationFunction, startInterval time.Duration, totalDuration time.Duration) chan interface{} {
	// Results channel will contain the results of all requests
	results := make(chan interface{})
	// Done channel will stop the main loop that send requests
	done := make(chan bool)

	var tickedFunctionWaitingGroup sync.WaitGroup
	tickedFunctionWaitingGroup.Add(1)

	// Call the ticked function in a goroutine to control its execution time
	go func() {
		defer func() {
			tickedFunctionWaitingGroup.Done()
		}()
		tickedFunctionRunner(tickedFunction, timeFunction, startInterval, done, results)

	}()

	// Wait until the end of the timer
	time.Sleep(totalDuration)

	// Stop the ticker
	done <- true

	// Wait for the last request before closing the results
	go func() {
		tickedFunctionWaitingGroup.Wait()
		close(results)
	}()

	return results
}

// tickedFunctionRunner creates a ticker and calls a function at every tick and vary the tick duration at each tick
// depending on the given variation function. The function will stop only when the done channel will receive a value.
// It returns a channel composed of every result of the tick function.
func tickedFunctionRunner(tickedFunction TickedFunction, timeFunction tickVariationFunction, startInterval time.Duration, done chan bool, results chan interface{}) {
	var requestsWaitingGroup sync.WaitGroup

	ticker := time.NewTicker(startInterval)
	defer ticker.Stop()

	tickIndex := 0
	previousTickDuration := startInterval

	startTime := time.Now()

	func() {
		for {
			// select allows to listen to channels
			select {
			case <-done:
				// Stop the loop
				return
			case t := <-ticker.C:
				// New tick ! Run function
				requestsWaitingGroup.Add(1)

				nextTickDuration := timeFunction(t.Sub(startTime))
				go func() {
					defer requestsWaitingGroup.Done()
					results <- tickedFunction(TickFunctionParameter{
						TickCount: tickIndex,
						Tick:      t,
						NextTick:  nextTickDuration,
					})
				}()

				// Prevent updating the ticker if not needed to make it more regular when using constants
				if nextTickDuration != previousTickDuration {
					previousTickDuration = nextTickDuration
					// Set the new ticker time
					ticker.Reset(nextTickDuration)
				}

			}
			tickIndex += 1
		}
	}()
	ticker.Stop()

	// Wait for the last requests before exiting
	requestsWaitingGroup.Wait()

}

// Linear calls a function multiple times during an interval starting at rpsStart call be second and ending at rpsEnd call
// per second.
func Linear(tickedFunction TickedFunction, rpsStart, rpsEnd int, totalDuration time.Duration) chan interface{} {
	// Calculate the coefficient
	startInterval := time.Second / time.Duration(rpsStart)

	linearTimeFunction := makeLinearVariationFunction(rpsStart, rpsEnd, totalDuration)
	return tickFunction(tickedFunction, linearTimeFunction, startInterval, totalDuration)
}
