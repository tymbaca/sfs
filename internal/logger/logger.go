// TEMPORARY PACKAGE
package logger

import (
	"fmt"
	"strings"
)

var Enabled = true

func Log(a ...any) (n int, err error) {
	if Enabled {
		return fmt.Println(a...)
	}

	return 0, nil
}

func Logf(format string, a ...any) (n int, err error) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if Enabled {
		return fmt.Printf(format, a...)
	}

	return 0, nil
}
