package graphql

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

// Output satisfies HTTP GraphQL request body.
type Output struct {
	Query         string `json:"query"`
	OperationName string `json:"operationName"`
	Variables     string `json:"variables"`
}

// Reset resets the default buffer for reusability.
func (e *Encoder) Reset() error {
	if _, ok := e.buf.(*bytes.Buffer); !ok {
		return errors.New("overridden io.Writer cannot be reset")
	}

	e.buf = bytes.NewBuffer(nil)
	return nil
}

// Marshal ...
func (e *Encoder) Marshal() ([]byte, error) {
	p, ok := e.buf.(Publisher)
	if !ok {
		return nil, errors.New("method not compatible with overridden io.Writer")
	}

	if len(e.objects) == 0 {
		return nil, errors.New("no items provided")
	}

	err := e.marshal(e.objects[0].inputSource, e.objects[0].queryName, e.objects[0].alias)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal input")
	}

	if ok {
		m, err := json.Marshal(&Output{
			Query: string(p.Bytes()),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request into json")
		}
		return m, nil
	}

	return nil, nil
}

// Query ...
func (e *Encoder) Query() ([]byte, error) {
	p, ok := e.buf.(Publisher)
	if !ok {
		return nil, errors.New("method not compatible with overridden io.Writer")
	}

	if len(e.objects) == 0 {
		return nil, errors.New("no items provided")
	}

	err := e.marshal(e.objects[0].inputSource, e.objects[0].queryName, e.objects[0].alias)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal input")
	}

	if ok {
		return p.Bytes(), nil
	}

	return nil, nil
}
