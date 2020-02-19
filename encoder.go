package graphql

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Encoder ...
type Encoder struct {
	objects     []encoderItem
	config      *config
	buf         io.Writer
	bufOverride bool
	log         *logrus.Entry
	logger      *logrus.Logger
}

type encoderItem struct {
	alias        string
	queryName    string
	outputSource interface{}
	inputSource  interface{}
}

// NewEncoder returns a new Encoder object.
func NewEncoder(requestType Type, prefix, indent string, opts ...Opt) (*Encoder, error) {
	if !requestType.isValid() {
		return nil, errors.New("invalid request type")
	}

	c, err := newConfig(requestType, prefix, indent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize config")
	}

	e := &Encoder{
		config:  c,
		objects: []encoderItem{},
		buf:     bytes.NewBuffer(nil),
	}

	e.initLog()

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

func (e *Encoder) initLog() {
	e.logger = logrus.New()
	e.logger.SetLevel(logrus.InfoLevel)
	e.logger.SetOutput(os.Stdout)
	e.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	e.logger.SetOutput(ioutil.Discard)
	e.log = logrus.NewEntry(e.logger)
}

// AddItem ...
func (e *Encoder) AddItem(queryName, alias string, variables interface{}, output interface{}) {
	e.objects = append(e.objects, encoderItem{
		alias:        alias,
		queryName:    queryName,
		inputSource:  variables,
		outputSource: output,
	})
}

// GetWriter returns the writer.
func (e *Encoder) GetWriter() io.Writer {
	return e.buf
}

// WriteString satisfies io.StringWriter interface.
func (e *Encoder) writeString(s string) error {
	// check if writer is StringWriter
	if sw, ok := e.buf.(io.StringWriter); ok {
		_, err := sw.WriteString(s)
		if err != nil {
			return errors.Wrap(err, "failed to write string")
		}
		return nil
	}

	// Handle not string optimized buffers
	_, err := fmt.Fprint(e.buf, s)
	if err != nil {
		return errors.Wrap(err, "failed to write string")
	}

	return nil
}

func (e *Encoder) writeOpenBracket() error {
	if e.config.indent != "" {
		err := e.writeString(e.config.inlineSpace + "{")
		if err != nil {
			return err
		}

		err = e.writeString("\n")
		if err != nil {
			return err
		}
		return nil
	}

	err := e.writeString("{")
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) writeCloseBracket(level int) error {
	if e.config.indent != "" {
		err := e.writeString(fmt.Sprintf("%s%s%s", e.config.prefix, e.getIndent(level), "}"))
		if err != nil {
			return err
		}

		err = e.writeString("\n")
		if err != nil {
			return err
		}

		return nil
	}

	err := e.writeString("}")
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) getIndent(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += e.config.indent
	}

	return s
}
