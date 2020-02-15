package main

import (
	"fmt"
	gql "graphql"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Foo ...
type Foo struct {
	GQLType gql.Type
	Bar     string `abc:"bar,out"`
	Foo     string `abc:"foo,out"`
}

func main() {
	f := Foo{
		GQLType: gql.TypeQuery(),
		Foo:     "test_foo",
		Bar:     "test_bar",
	}

	marshalled, err := gql.MarshalIndent(&f, "", "", gql.TagNameInject("abc"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to marshal gql"))
	}

	fmt.Println(string(marshalled))
}
