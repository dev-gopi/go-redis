package storage

type JSONValue struct {
	Raw    []byte
	Parsed map[string]any
}
