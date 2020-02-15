package main

import (
	"fmt"
	"graphql-client"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Foo ...
type Foo struct {
	Bar string `gql:"bar,out"`
	Foo string `gql:"foo,out"`
}

func main() {
	f := Foo{
		Foo: "test_foo",
		Bar: "test_bar",
	}

	marshalled, err := graphql.MarshalIndent(&f, "", "")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to marshal gql"))
	}

	fmt.Println(string(marshalled))
}
