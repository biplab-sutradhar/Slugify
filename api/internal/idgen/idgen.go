package idgen

import (
	"errors"
	"fmt"
)

var counter int64

func GenerateShortCode() (string, error) {
	if counter == int64(^uint(0)>>1) { // Check if the counter is at max int64 value
		return "", errors.New("counter overflow")
	}
	counter++
	return fmt.Sprintf("%d", counter), nil
}
