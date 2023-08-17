package jsonstream

import (
	"github.com/goccy/go-json"
	"io"
)

var GoJSONConfig = &Config{
	NewEncoder: NewGoJSONEncoder,
	NewDecoder: NewGoJSONDecoder,
}

func NewGoJSONEncoder(out io.Writer) JSONEncoder {
	return json.NewEncoder(out)
}

type goJSONDecoder struct {
	*json.Decoder
}

func NewGoJSONDecoder(in io.Reader) JSONDecoder {
	return &goJSONDecoder{json.NewDecoder(in)}
}

func (d *goJSONDecoder) RuneToken() (token rune, others any, err error) {
	t, err := d.Decoder.Token()
	if err != nil {
		return 0, nil, err
	}
	if d, ok := t.(json.Delim); ok {
		return rune(d), nil, nil
	}
	return 0, t, nil
}
