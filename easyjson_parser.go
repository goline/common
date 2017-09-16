package utils

import (
	j "github.com/mailru/easyjson"
)

type EasyJsonParser struct {}

func (p *EasyJsonParser) ContentType() string {
	return "application/json"
}

func (p *EasyJsonParser) Decode(data []byte, v interface{}) error {
	return j.Unmarshal(data, v.(j.Unmarshaler))
}

func (p *EasyJsonParser) Encode(v interface{}) ([]byte, error) {
	return j.Marshal(v.(j.Marshaler))
}