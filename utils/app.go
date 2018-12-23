package utils

import (
	"fmt"
	"time"
)

func RunTime(now time.Time) {
	fmt.Printf("\nApplication took %dns to complete\n", time.Since(now))
}
