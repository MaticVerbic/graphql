package main

import (
	"fmt"
	gql "graphql"
	"time"
)

// Foo ...
type Foo struct {
	NameField string
	Foo       string `abc:"foo"`
	Bar       Bar    `abc:"bar"`
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
	}

	//b := bytes.NewBuffer(nil)
	start := time.Now()
	enc, err := gql.NewEncoder(gql.TypeQuery(), "", "  ",
		gql.TagNameOpt("abc"),
		gql.NameFieldOpt("NameField"),
		//gql.OverrideWriterOpt(os.Stdout),
		//gql.LogLevelOpt(logrus.DebugLevel),
		//gql.LogOutputOpt(os.Stdout),
	)
	if err != nil {
		panic(err)
	}

	enc.AddItem("", "", &f, &f)

	marshalled, err := enc.Marshal()
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start).Nanoseconds()
	//fmt.Println(b.String())
	fmt.Println(string(marshalled))
	fmt.Printf("elapsed: %dns", elapsed)
}
