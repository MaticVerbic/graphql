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

	e.writeString(e.config.prefix + e.config.requestType.String())
	e.writeOpenBracket()
	switch t.Kind() {
	case reflect.Struct:
		e.writeString(e.config.prefix + e.config.indent + e.getName(source))
		e.writeOpenBracket()

		e.handleStruct(source, 2)

		e.writeCloseBracket(1)

	case reflect.Map:
		break
	default:
		return errors.New("invalid source type")
	}

	e.writeCloseBracket(0)
	return nil
}

func (e *Encoder) handleStruct(s interface{}, level int) {
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
			e.writeString(e.config.prefix + e.getIndent(level) + tag)
			e.writeOpenBracket()

			// recursively handle child structs
			e.handleStruct(v.Field(i).Addr().Interface(), level+1)

			// close a new recursion level
			e.writeCloseBracket(level)
			continue
		}

		// write a simple field name
		e.writeString(e.config.prefix + e.getIndent(level) + tag + "\n")
	}

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
