package main

import (
	"fmt"
	gql "graphql"
)

// Foo ...
type Foo struct {
	NameField string
	Bar       string `abc:"bar,out"`
	Foo       string `abc:"foo,out"`
}

func main() {
	f := Foo{
		NameField: "test",
		Foo:       "test_foo",
		Bar:       "test_bar",
	}

	enc, err := gql.NewEncoder(gql.TypeQuery(), "", "  ",
		gql.TagNameOpt("abc"),
		gql.NameFieldOpt("NameField"),
	)
	if err != nil {
		panic(err)
	}

	enc.AddItem("", &f, &f)

	marshalled, err := enc.MarshalIndent()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshalled))

	//fmt.Println(string(marshalled))
}
