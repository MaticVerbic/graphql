package graphql

import (
	"reflect"

	"github.com/pkg/errors"
)

func (e *Encoder) marshal(source interface{}) error {
	if source == nil {
		return errors.New("source is nil interface")
	}

	v := reflect.Indirect(reflect.ValueOf(source))
	t := reflect.TypeOf(v)

	err := e.writeString(e.config.prefix + e.config.requestType.String())
	if err != nil {
		return err
	}

	err = e.writeOpenBracket()
	if err != nil {
		return err
	}

	switch t.Kind() {
	case reflect.Struct:
		err = e.writeString(e.config.prefix + e.config.indent + e.getName(source))
		if err != nil {
			return err
		}

		err = e.writeOpenBracket()
		if err != nil {
			return err
		}

		err = e.handleStruct(source, 2)
		if err != nil {
			return err
		}

		err = e.writeCloseBracket(1)
		if err != nil {
			return err
		}

	case reflect.Map:
		break
	default:
		return errors.New("invalid source type")
	}

	err = e.writeCloseBracket(0)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) handleStruct(s interface{}, level int) error {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := reflect.TypeOf(s)
	t = t.Elem()

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
			err := e.writeString(e.config.prefix + e.getIndent(level) + tag)
			if err != nil {
				return err
			}

			err = e.writeOpenBracket()
			if err != nil {
				return err
			}

			// recursively handle child structs
			err = e.handleStruct(v.Field(i).Addr().Interface(), level+1)
			if err != nil {
				return err
			}

			// close a new recursion level
			err = e.writeCloseBracket(level)
			if err != nil {
				return err
			}

			continue
		}

		// write a simple field name
		err := e.writeString(e.config.prefix + e.getIndent(level) + tag + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) getName(s interface{}) string {
	v := reflect.ValueOf(s).Elem()
	gt := v.FieldByName(e.config.nameField)
	if !gt.IsValid() && reflect.TypeOf(gt).Kind() != reflect.String {
		s := reflect.TypeOf(s).Elem().Name()
		return s
	}

	gtt := gt.Interface().(string)

	return gtt
}
