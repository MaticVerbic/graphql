package graphql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func (e *Encoder) marshal(source interface{}) ([]byte, error) {
	if source == nil {
		return nil, errors.New("source is nil interface")
	}

	v := reflect.Indirect(reflect.ValueOf(source))
	t := reflect.TypeOf(v)

	querySkeleton := `%s {
  %s {%s
  }
}`

	var queryBody string
	var queryName string
	var err error
	switch t.Kind() {
	case reflect.Struct:
		queryName = e.getName(source)

		queryBody, err = e.handleStruct(source)
		if err != nil {
			return nil, errors.Wrap(err, "failed to handle a struct")
		}
	case reflect.Map:
		break
	default:
		return nil, errors.New("invalid source type")
	}

	return []byte(fmt.Sprintf(querySkeleton, e.config.requestType, queryName, queryBody)), nil
}

func (e *Encoder) handleStruct(s interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := reflect.TypeOf(s)
	t = t.Elem()

	str := ""
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		tag := ft.Tag.Get(e.config.tagname)

		// @todo: better handling for empty tags.
		if tag == "" {
			continue
		}

		spl := strings.Split(tag, ",")
		if len(spl) != 2 {
			return "", errors.New("invalid separator count")
		}

		if spl[1] != "out" {
			continue
		}

		str += "\n" + fmt.Sprintf(`%8s`, spl[0])
	}

	return str, nil
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
