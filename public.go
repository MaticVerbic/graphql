package graphql

import "github.com/pkg/errors"

// Encoder ...
type Encoder struct {
	objects []encoderItem
	config  *config
}

type encoderItem struct {
	alias        string
	outputSource interface{}
	inputSource  interface{}
}

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

// NewEncoder returns a new Encoder object.
func NewEncoder(requestType Type, prefix, indent string, opts ...Opt) (*Encoder, error) {
	const (
		defaultTag = "gql"
	)

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
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

// AddItem ...
func (e *Encoder) AddItem(alias string, variables interface{}, output interface{}) {
	e.objects = append(e.objects, encoderItem{
		alias:        alias,
		inputSource:  variables,
		outputSource: output,
	})
}

// MarshalIndent ...
func (e *Encoder) MarshalIndent() ([]byte, error) {
	return e.marshal(e.objects[0].inputSource) // @todo: marshalIndent.
}

// Marshal ...
func (e *Encoder) Marshal(source interface{}) ([]byte, error) {
	return e.marshal(e.objects[0].inputSource)
}
