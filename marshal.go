package graphql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func marshal(source interface{}, c *config) ([]byte, error) {
	if source == nil {
		return nil, errors.New("source is nil interface")
	}

	v := reflect.Indirect(reflect.ValueOf(source))
	t := reflect.TypeOf(v)

	querySkeleton := `
%s {
  %s {%s
  }
}`

	var queryBody string
	var err error
	switch t.Kind() {
	case reflect.Struct:
		queryBody, err = handleStruct(source, c.tagname, c)
		if err != nil {
			return nil, errors.Wrap(err, "failed to handle a struct")
		}
	case reflect.Map:
		break
	default:
		return nil, errors.New("invalid source type")
	}

	return []byte(fmt.Sprintf(querySkeleton, c.typ, c.requestName, queryBody)), nil
}

func handleStruct(s interface{}, tagname string, c *config) (string, error) {
	const gqlType = "GQLType"

	v := reflect.Indirect(reflect.ValueOf(s))
	t := reflect.TypeOf(s)
	t = t.Elem()

	str := ""
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		if ft.Name == gqlType {
			gt := Type(v.FieldByName(gqlType).String())
			if !gt.isValid() {
				return "", fmt.Errorf("invalid value for %q", gqlType)
			}

			c.SetType(gt)
		}

		tag := ft.Tag.Get(tagname)

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
