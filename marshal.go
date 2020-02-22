package graphql

import (
	"reflect"

	"github.com/pkg/errors"
)

// ErrGeneral is a general error returned by marshaler
const ErrGeneral = "failed to marshal interface"

func (e *Encoder) marshal(source interface{}, name, alias string) error {
	if source == nil {
		return errors.New("source is nil interface")
	}

	v := reflect.Indirect(reflect.ValueOf(source))
	t := reflect.TypeOf(v.Interface())

	if err := e.writeObjectHeader(0, e.config.requestType.String()); err != nil {
		return errors.Wrap(err, "failed to write object header")
	}

	if err := e.writeObjectHeader(1, name); err != nil {
		return errors.Wrap(err, "failed to write object header")
	}

	switch t.Kind() {
	case reflect.Struct:
		if err := e.handleStruct(source, 2); err != nil {
			return errors.Wrap(err, ErrGeneral)
		}

	case reflect.Map:
		if err := e.handleMap(source, 2); err != nil {
			return errors.Wrap(err, ErrGeneral)
		}
	default:
		return errors.New("invalid source type")
	}

	if err := e.writeCloseBracket(1); err != nil {
		return errors.Wrap(err, ErrGeneral)
	}

	if err := e.writeCloseBracket(0); err != nil {
		return errors.Wrap(err, ErrGeneral)
	}

	return nil
}

func (e *Encoder) handleMap(m interface{}, level int) error {
	var ma map[string]interface{}
	if ca, ok := m.(*map[string]interface{}); ok {
		ma = *ca
	} else if ca, ok := m.(map[string]interface{}); ok {
		ma = ca
	} else {
		return errors.New("invalid map type")
	}

	var err error
	var inlineCount int
	for key, value := range ma {
		v := reflect.Indirect(reflect.ValueOf(value))
		inlineCount = 0

		switch v.Kind() {
		case reflect.String:
			if isValid(v) {
				if err = e.writeObjectHeader(level, key); err != nil {
					return err
				}

				inlineCount, err = e.writeItem(inlineCount, level+1, value.(string))
				if err != nil {
					return err
				}

				if err = e.writeCloseBracket(level); err != nil {
					return err
				}
			} else {
				inlineCount, err = e.writeItem(inlineCount, level, key)
				if err != nil {
					return err
				}
			}
		case reflect.Slice:
			if isValid(v) {
				if err = e.writeObjectHeader(level, key); err != nil {
					return err
				}

				for _, item := range v.Interface().([]string) {
					inlineCount, err = e.writeItem(inlineCount, level+1, item)
					if err != nil {
						return err
					}
				}

				if err = e.writeCloseBracket(level); err != nil {
					return err
				}
			} else {
				inlineCount, err = e.writeItem(inlineCount, level, key)
				if err != nil {
					return err
				}
			}
		case reflect.Struct:
			if err = e.writeObjectHeader(level, key); err != nil {
				return err
			}

			if err = e.handleStruct(value, level+1); err != nil {
				return err
			}

			if err = e.writeCloseBracket(level); err != nil {
				return err
			}
		case reflect.Map:
			if isValid(v) {
				if err = e.writeObjectHeader(level, key); err != nil {
					return err
				}

				if err = e.handleMap(value, level+1); err != nil {
					return err
				}

				if err = e.writeCloseBracket(level); err != nil {
					return err
				}
			} else {
				inlineCount, err = e.writeItem(inlineCount, level, key)
				if err != nil {
					return err
				}
			}
		default:
			continue
		}
	}

	return nil
}

func (e *Encoder) handleStruct(s interface{}, level int) error {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := reflect.TypeOf(v.Interface())

	var err error
	inlineCount := 0
	for i := 0; i < v.NumField(); i++ {
		fv := reflect.Indirect(reflect.ValueOf(v.Field(i).Interface()))
		ft := t.Field(i)
		tag := ft.Tag.Get(e.config.tagname)

		// @todo: better handling for empty tags.
		if tag == "" {
			continue
		}

		// custom handling of structs, maps and arrays
		switch v.Field(i).Kind() {
		case reflect.Struct:
			// set up a new recursion level
			if err = e.writeObjectHeader(level, tag); err != nil {
				return errors.Wrap(err, ErrGeneral)
			}

			// recursively handle child structs
			if err = e.handleStruct(fv.Interface(), level+1); err != nil {
				return errors.Wrapf(err, "failed to handle struct %q", tag)
			}

			// close a new recursion level
			if err = e.writeCloseBracket(level); err != nil {
				return errors.Wrap(err, ErrGeneral)
			}

			continue
		case reflect.Map:
			if isValid(v.Field(i)) {
				if err = e.writeObjectHeader(level, tag); err != nil {
					return err
				}

				if err = e.handleMap(v.Field(i).Interface(), level+1); err != nil {
					return err
				}

				if err = e.writeCloseBracket(level); err != nil {
					return err
				}
			}
		default:
			inlineCount, err = e.writeItem(inlineCount, level, tag)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Encoder) getName(s interface{}) string {
	v := reflect.ValueOf(s).Elem()
	gt := v.FieldByName(e.config.nameField)

	if (!gt.IsValid() && reflect.TypeOf(gt).Kind() != reflect.String) ||
		gt.Interface().(string) == "" {
		s := reflect.TypeOf(s).Elem().Name()
		return s
	}

	return gt.Interface().(string)
}

func isValid(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if !v.IsNil() {
			return true
		}
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		if !(v.Len() == 0) {
			return true
		}
	}

	return false
}
