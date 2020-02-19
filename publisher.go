package graphql

// Publisher interface allows returning of
// []byte in Marshal and MarshalIndent methods
type Publisher interface {
	Bytes() []byte
}
