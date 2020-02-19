package graphql

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

// Opt is an option func.
type Opt func(e *Encoder)

// TagNameOpt injects a custom tag
func TagNameOpt(tagname string) Opt {
	return func(e *Encoder) {
		e.config.tagname = tagname
	}
}

// NameFieldOpt injects a custom tag
func NameFieldOpt(nameField string) Opt {
	return func(e *Encoder) {
		e.config.nameField = nameField
	}
}

// OverrideWriterOpt overrides a new bytes.Buffer in Encoder.
func OverrideWriterOpt(buf io.Writer) Opt {
	return func(e *Encoder) {
		e.buf = buf
		e.bufOverride = true
	}
}

// InlineSpaceOpt overrides a single space separator
func InlineSpaceOpt(s string) Opt {
	return func(e *Encoder) {
		if s != "" && !strings.ContainsAny(s, " ") {
			return
		}

		e.config.inlineSpace = s
	}
}

// OverrideLogOpt allows for custom loggers
func OverrideLogOpt(log *logrus.Entry) Opt {
	return func(e *Encoder) {
		if log == nil {
			return
		}

		e.log = log
	}
}

// LogLevelOpt allows for custom loggers
func LogLevelOpt(level logrus.Level) Opt {
	return func(e *Encoder) {
		e.logger.SetLevel(level)
		e.log = logrus.NewEntry(e.logger)
	}
}

// LogOutputOpt allows for custom loggers
func LogOutputOpt(w io.Writer) Opt {
	return func(e *Encoder) {
		if w == nil {
			return
		}

		e.logger.SetOutput(w)
		e.log = logrus.NewEntry(e.logger)
	}
}
