package utils

import (
	"errors"
	"time"
)

// Retry retries the callback function a specified number of times
func Retry(times int, callback func(int) error, sleepMilliseconds int, when func(error) bool) error {
	attempts := 0
	var backoff []int

	// Check if times is an array (backoff values)
	if times > 0 {
		backoff = make([]int, times)
	}

	// Loop for the number of times specified
	for times > 0 {
		attempts++
		times--

		err := callback(attempts)
		if err == nil {
			return nil
		}

		// Check if retries are exhausted or the condition fails
		if times < 1 || (when != nil && !when(err)) {
			return err
		}

		sleepMillisecondsVal := sleepMilliseconds
		if len(backoff) > attempts-1 {
			sleepMillisecondsVal = backoff[attempts-1]
		}

		if sleepMillisecondsVal > 0 {
			time.Sleep(time.Duration(sleepMillisecondsVal) * time.Millisecond)
		}
	}

	return errors.New("max retries exceeded")
}
