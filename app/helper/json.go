package helper

import (
	"github.com/goccy/go-json"
)

func JSONMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func JSONMarshalNoEscape(v interface{}) ([]byte, error) {
	return json.MarshalNoEscape(v)
}

func JSONMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
