package helper

import (
	"io"
	"os"
	"strings"
)

// errToString converts an error to a string, returning "nil" if error is nil
func ErrToString(err error) string {
	if err == nil {
		return "nil"
	}
	return err.Error()
}

func FiberLoggerOutput(str string) io.Writer {
	var output io.Writer
	switch strings.ToLower(str) {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "file":
		// In a real implementation, you would open a file here
		// For now, default to stdout
		output = os.Stdout
	default:
		output = os.Stdout
	}
	return output
}
