package common

import (
	"fmt"
	"os"
)

func ExitWithError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(2)
}
