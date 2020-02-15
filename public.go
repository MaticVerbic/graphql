package graphql

import "github.com/pkg/errors"

// MarshalIndent ...
func MarshalIndent(source interface{}, prefix, indent string) ([]byte, error) {
	c := config{prefix, indent}
	if err := c.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate opt")
	}

	return marshal(source, c) // @todo: marshalIndent.
}

// Marshal ...
func Marshal(source interface{}) ([]byte, error) {
	return marshal(source, config{})
}
