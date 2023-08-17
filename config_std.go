package jsonstream

import (
	"encoding/json"
	"io"
)

var StdConfig = &Config{
	NewEncoder: NewStdJSONEncoder,
	NewDecoder: NewStdJSONDecoder,
}

func NewStdJSONEncoder(out io.Writer) JSONEncoder {
	return json.NewEncoder(out)
}

type stdJSONDecoder struct {
	*json.Decoder
}

func NewStdJSONDecoder(in io.Reader) JSONDecoder {
	return &stdJSONDecoder{json.NewDecoder(in)}
}

func (d *stdJSONDecoder) RuneToken() (token rune, others any, err error) {
	t, err := d.Decoder.Token()
	if err != nil {
		return 0, nil, err
	}
	if d, ok := t.(json.Delim); ok {
		return rune(d), nil, nil
	}
	return 0, t, nil
}
