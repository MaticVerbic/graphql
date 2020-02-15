package graphql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type config struct{ prefix, indent string }

func (c config) validate() error {
	if c.indent != "" && !strings.ContainsAny(c.indent, " ") {
		return errors.New("non whitespace char in 'indent' arg")
	}

	if c.prefix != "" && !strings.ContainsAny(c.prefix, " ") {
		return errors.New("non whitespace char in 'prefix' arg")
	}

	return nil
}

func marshal(source interface{}, c config) ([]byte, error) {
	if source == nil {
		return nil, errors.New("source is nil interface")
	}

	v := reflect.Indirect(reflect.ValueOf(source))
	t := reflect.TypeOf(v)

	switch t.Kind() {
	case reflect.Struct:
		return handleStruct(source, "gql")
	case reflect.Map:
		break
	default:
		return nil, errors.New("invalid source type")
	}

	return nil, nil
}

func handleStruct(s interface{}, tagname string) ([]byte, error) {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := reflect.TypeOf(s)

	t = t.Elem()

	querySkeleton := `
query {
  %s%s {%s
  }
}`

	args := ""
	str := ""
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)

		tag := ft.Tag.Get(tagname)
		spl := strings.Split(tag, ",")
		if len(spl) != 2 {
			return nil, errors.New("invalid separator count")
		}

		if spl[1] != "out" {
			continue
		}

		str += "\n" + fmt.Sprintf(`%8s`, spl[0])
	}

	return []byte(fmt.Sprintf(querySkeleton, strings.ToLower(t.Name()), args, str)), nil
}
