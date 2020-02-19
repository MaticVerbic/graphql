package graphql

import (
	"strings"

	"github.com/pkg/errors"
)

type config struct {
	tagname     string
	nameField   string
	inlineSpace string
	prefix      string
	indent      string
	requestType Type
}

func newConfig(requestType Type, prefix, indent string) (*config, error) {
	const (
		defaultTag       = "gql"
		defaultNameField = "GQLName"
	)

	c := &config{
		tagname:     defaultTag,
		nameField:   defaultNameField,
		inlineSpace: " ",
		prefix:      prefix,
		indent:      indent,
		requestType: requestType,
	}

	if err := c.validate(); err != nil {
		return nil, err
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
