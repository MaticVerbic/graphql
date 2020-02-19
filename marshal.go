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
	t := reflect.TypeOf(v)

	if err := e.writeString(e.config.prefix + e.config.requestType.String()); err != nil {
		return errors.Wrap(err, "failed to marshal request type")
	}

	if err := e.writeOpenBracket(); err != nil {
		return errors.Wrap(err, ErrGeneral)
	}

	if name != "" || t.Kind() == reflect.Struct {
		if err := e.writeString(e.config.prefix + e.config.indent); err != nil {
			return errors.Wrap(err, "failed to marshal request name indentation")
		}

		if alias != "" {
			if err := e.writeString(alias + ": "); err != nil {
				return errors.Wrapf(err, "failed to marshal alias %q", alias)
			}
		}

		if name == "" {
			if err := e.writeString(e.getName(source)); err != nil {
				return errors.Wrapf(err, "failed to marshal name %q", e.getName(source))
			}
		} else {
			if err := e.writeString(name); err != nil {
				return errors.Wrapf(err, "failed to marshal name %q", name)
			}
		}

		if err := e.writeOpenBracket(); err != nil {
			return errors.Wrap(err, ErrGeneral)
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		if err := e.handleStruct(source, 2); err != nil {
			return errors.Wrap(err, ErrGeneral)
		}

	case reflect.Map:
		break
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

func (e *Encoder) handleStruct(s interface{}, level int) error {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := reflect.TypeOf(s)
	t = t.Elem()

	inlineCount := 0

	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		tag := ft.Tag.Get(e.config.tagname)

		// @todo: better handling for empty tags.
		if tag == "" {
			continue
		}

		// custom handling of structs, maps and arrays
		switch ft.Type.Kind() {
		case reflect.Struct:
			// set up a new recursion level
			if e.config.indent != "" {
				if err := e.writeString(e.config.prefix + e.getIndent(level) + tag); err != nil {
					return errors.Wrapf(err, "failed to handle struct %q", tag)
				}
			} else {
				if inlineCount > 0 {
					if err := e.writeString(e.config.inlineSpace); err != nil {
						return errors.Wrap(err, ErrGeneral)
					}
				}

				if err := e.writeString(tag); err != nil {
					return errors.Wrapf(err, "failed to handle struct %q", tag)
				}

				inlineCount++
			}

			if err := e.writeOpenBracket(); err != nil {
				return errors.Wrap(err, ErrGeneral)
			}

			// recursively handle child structs
			if err := e.handleStruct(v.Field(i).Addr().Interface(), level+1); err != nil {
				return errors.Wrapf(err, "failed to handle struct %q", tag)
			}

			// close a new recursion level
			if err := e.writeCloseBracket(level); err != nil {
				return errors.Wrap(err, ErrGeneral)
			}

			continue
		case reflect.Map:
			continue
		default:
			if e.config.indent != "" {
				if err := e.writeString(e.config.prefix + e.getIndent(level) + tag + "\n"); err != nil {
					return errors.Wrapf(err, "failed to write field name %q", tag)
				}
			} else {
				if inlineCount > 0 {
					if err := e.writeString(e.config.inlineSpace); err != nil {
						return errors.Wrap(err, ErrGeneral)
					}
				}

				if err := e.writeString(tag); err != nil {
					return errors.Wrapf(err, "failed to write field name %q", tag)
				}

				inlineCount++
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
