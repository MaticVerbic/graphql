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

type config struct {
	tagname     string
	prefix      string
	indent      string
	requestName string
	typ         Type
}

func (c *config) SetType(t Type) { c.typ = t }

func newConfig(prefix, indent string, opts ...Opt) (*config, error) {
	const (
		defaultTag         = "gql"
		defaultRequestName = "queryObject"
	)

	c := &config{
		tagname:     defaultTag,
		prefix:      prefix,
		indent:      indent,
		requestName: defaultRequestName,
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
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
