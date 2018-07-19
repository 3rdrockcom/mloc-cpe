package log

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/epointpayment/mloc-cpe/app/config"

	"github.com/juju/errors"
)

var pathPrefix = strings.Join([]string{os.Getenv("GOPATH"), "src"}, string(os.PathSeparator)) + string(os.PathSeparator)

type Stack []StackFrame

type StackFrame struct {
	Filename   string
	LineNumber int64
	Context    string
	Underlying string
}

func (s Stack) String() (output string) {
	for i, frame := range s {
		o := fmt.Sprintf("#%v: %s:%d\n", i, frame.Filename, frame.LineNumber)

		if frame.Context != "" {
			o += fmt.Sprintf("\t%s: '%s'\n", "Context", frame.Context)
		}

		if frame.Underlying != "" {
			o += fmt.Sprintf("\t%s: %s\n", "Cause", frame.Underlying)
		}

		output += o
	}

	return
}

func StackTrace(Error error) (s Stack) {
	Err, ok := Error.(*errors.Err)
	if !ok {
		return
	}

	for _, entry := range Err.StackTrace() {
		if entry == "" || (Err.Underlying() != nil && entry == Err.Underlying().Error()) {
			continue
		}

		f := StackFrame{}

		parts := strings.Split(entry, ":")
		switch {
		case len(parts) == 3:
			if n, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				f.Filename = prependPath(parts[0])
				f.LineNumber = n
				f.Context = strings.TrimSpace(parts[2])
			}
		case len(parts) == 2:
			if n, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				f.Filename = prependPath(parts[0])
				f.LineNumber = n
			}
		default:
			// Fallback to avoid erroring out here if no location is found
			f.Filename = entry
		}

		s = append(s, f)
	}

	if len(s) > 0 && Err.Underlying() != nil {
		s[len(s)-1].Underlying = strings.TrimSpace(Err.Underlying().Error())
	}

	return s
}

func prependPath(filename string) string {
	if environment == config.EnvDevelopment {
		return pathPrefix + strings.TrimPrefix(filename, pathPrefix)
	}

	return filename
}
