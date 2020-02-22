package main

import (
	"fmt"
	gql "graphql"
	"time"
)

// Foo ...
type Foo struct {
	NameField string
	Foo       string                 `abc:"foo"`
	Bar       Bar                    `abc:"bar"`
	Map       map[string]interface{} `abc:"map"`
}

// Bar ...
type Bar struct {
	Fb  string `abc:"fb"`
	Baz Baz    `abc:"baz"`
}

// Baz ...
type Baz struct {
	Foobar string `abc:"foobar"`
}

func main() {
	f := Foo{
		NameField: "test",
		Foo:       "test_foo",
		Bar: Bar{
			Fb: "fb",
			Baz: Baz{
				Foobar: "foobar",
			},
		},
		Map: map[string]interface{}{"someKey": "someValue"},
	}

	m := map[string]interface{}{
		"struct":       &f,
		"sliceString":  &[]string{"foo", "bar", "baz"},
		"stringSingle": "foobar",
		"sliceEmpty":   []string{},
		"stringEmpty":  "",
		"map":          map[string]interface{}{"someKey": "someValue"},
	}

	start := time.Now()
	enc, err := gql.NewEncoder(gql.TypeQuery(), "", "  ",
		gql.TagNameOpt("abc"),
		gql.NameFieldOpt("NameField"),
	)
	if err != nil {
		panic(err)
	}

	if err := enc.AddItem("testNamedQuery", "", &f, &f); err != nil {
		panic(err)
	}

	marshalled, err := enc.Marshal()
	if err != nil {
		panic(err)
	}

	err = enc.Reset()
	if err != nil {
		panic(err)
	}

	query, err := enc.Query()
	if err != nil {
		panic(err)
	}

	enc, err = gql.NewEncoder(gql.TypeQuery(), "", "  ",
		gql.TagNameOpt("abc"),
		//gql.OverrideWriterOpt(os.Stdout),
		//gql.LogLevelOpt(logrus.DebugLevel),
		//gql.LogOutputOpt(os.Stdout),
	)
	if err != nil {
		panic(err)
	}

	if err := enc.AddItem("queryNameTest", "", &m, &m); err != nil {
		panic(err)
	}

	marshalledTwo, err := enc.Marshal()
	if err != nil {
		panic(err)
	}

	err = enc.Reset()
	if err != nil {
		panic(err)
	}

	queryTwo, err := enc.Query()
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start).Microseconds()
	//fmt.Println(b.String())
	fmt.Println("http.Request ready: \n" + string(marshalled))
	fmt.Println("\nclient ready: \n" + string(query))
	fmt.Println("http.Request ready: \n" + string(marshalledTwo))
	fmt.Println("\nclient ready: \n" + string(queryTwo))
	fmt.Printf("elapsed: %dms", elapsed)
}
