package graphql

import "github.com/pkg/errors"

// MarshalIndent ...
func MarshalIndent(source interface{}, prefix, indent string, opts ...Opt) ([]byte, error) {
	c, err := newConfig(prefix, indent, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new config")
	}

	return marshal(source, c) // @todo: marshalIndent.
}

// Marshal ...
func Marshal(source interface{}) ([]byte, error) {
	c, _ := newConfig("", "")
	return marshal(source, c)
}
