package migrations

import (
	"strings"

	"github.com/epointpayment/mloc-cpe/app/log"
)

// Verbose displays additional logging messages if enabled
var Verbose bool

// Logger is an interface so you can pass in your own
// logging implementation.
type Logger struct{}

// Printf is like fmt.Printf
func (l *Logger) Printf(format string, v ...interface{}) {
	log.Printf(strings.TrimRight("migrate: "+format, "\n"), v...)
}

// Verbose should return true when verbose logging output is wanted
func (l *Logger) Verbose() bool {
	return Verbose
}
