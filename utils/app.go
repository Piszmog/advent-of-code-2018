package utils

import (
	"fmt"
	"time"
)

func Runtime(now time.Time) {
	fmt.Printf("\nApplication took %dns to complete\n", time.Since(now))
}
