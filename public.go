package graphql

import "github.com/pkg/errors"

// Marshal ...
func (e *Encoder) Marshal() ([]byte, error) {
	p, ok := e.buf.(Publisher)
	if !ok {
		e.log.Warn("output writer does not implement Publisher interface, output will be nil")
	}

	err := e.marshal(e.objects[0].inputSource, e.objects[0].queryName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal input")
	}

	if ok {
		return p.Bytes(), nil
	}

	return nil, nil
}
