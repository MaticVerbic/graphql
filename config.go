package graphql

import (
	"strings"

	"github.com/pkg/errors"
)

// Opt is an option func.
type Opt func(c *config)

// TagNameInject injects a custom tag
func TagNameInject(tagname string) Opt {
	return func(c *config) {
		c.tagname = tagname
	}
}

type config struct{ tagname, prefix, indent string }

func newConfig(opts ...Opt) *config {
	const (
		defaultPrefix = ""
		defaultIndent = "  "
		defaultTag    = "gql"
	)

	c := &config{
		tagname: defaultTag,
		prefix:  defaultPrefix,
		indent:  defaultIndent,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *config) validate() error {
	if c.indent != "" && !strings.ContainsAny(c.indent, " ") {
		return errors.New("non whitespace char in 'indent' arg")
	}

	if c.prefix != "" && !strings.ContainsAny(c.prefix, " ") {
		return errors.New("non whitespace char in 'prefix' arg")
	}

	return nil
}
