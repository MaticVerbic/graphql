package graphql

// Type string.
type Type string

const (
	typeQuery    Type = "query"
	typeMutation Type = "mutation"
)

// TypeQuery method.
func TypeQuery() Type { return typeQuery }

// TypeMutation method.
func TypeMutation() Type { return typeMutation }

func (t Type) isValid() bool { return t == typeMutation || t == typeQuery }

func (t Type) String() string { return string(t) }
