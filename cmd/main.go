package main

import (
	"fmt"
	"graphql/graphql"
)

// Foo ...
type Foo struct {
	Bar string `gql:"bar,out"`
	Foo string `gql:"foo,out"`
}

func main() {
	marshalled, err := graphql.MarshalIndent(&Foo{Foo: "", Bar: ""}, "", "")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshalled))
}
