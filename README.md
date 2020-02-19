# graphql

Package that offers simple support for marshalling go datatypes into GraphQL requests.

## Usage

```go
// Foobar ...
type Foobar struct {
  GQLName   string
  Foo       string `gql:"foo"`
  Bar       string `gql:"bar"`
}


func main() {
  // Create a new output object.
  f := Foobar{
    GQLName: "test",
  }

  // Create a new graphql Encoder.
  // Encoder takes graphql.Type, prefix and indent as arguments.
  // If indent is an empty string, then the output will be a single line.
  enc, err := gql.NewEncoder(gql.TypeQuery(), "", "  ")
  if err != nil {
    panic(err)
  }

  // Add a  request.
  // GraphQL supports multiple requests per single http.Request.
  // If name not provided, then the encoder will revert to GQLName field,
  // if that is not found, it will revert to provided struct name, otherwise throw an error.
  enc.AddItem("", "", nil, &f)

  // Marshal the request
  marshalled, err := enc.Marshal()
  if err != nil {
    panic(err)
  }
```

### Opts

```go
enc, err := gql.NewEncoder(gql.TypeQuery(), "", "  ",
    gql.TagNameOpt("abc"),                                // overrides default tag used for parsing; default: gql
    gql.NameFieldOpt("NameField"),                        // overrides default field name used in query naming; default: GQLName
    gql.InlineSpaceOpt("   "),                            // overrides the spacing used between two adjacent elements on a single line; default: " "
    gql.OverrideWriterOpt(os.Stdout),                     // overrides default io.Writer; default: bytes.Buffer
    gql.OverrideLogOpt(logrus.NewEntry(logrus.New())),    // overrides default logrus.Entry
    gql.LogLevelOpt(logrus.WarnLevel),                    // overrides default logrus.Level; default: logrus.InfoLevel
    gql.LogOutputOpt(os.Stdout),                          // overrides default logrus output; default: ioutil.Discard
)
```
