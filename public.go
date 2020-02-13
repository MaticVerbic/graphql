package graphql

import "github.com/pkg/errors"

// MarshalIndent ...
func MarshalIndent(source interface{}, prefix, indent string) ([]byte, error) {
	o := opt{prefix, indent}
	if err := o.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate opt")
	}

	return marshal(source, opt{prefix, indent}) // @todo: marshalIndent.
}

// Marshal ...
func Marshal(source interface{}) ([]byte, error) {
	return marshal(source, opt{})
}
