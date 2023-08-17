package jsonstream

import (
	"io"
)

type (
	Config struct {
		NewEncoder func(out io.Writer) JSONEncoder
		NewDecoder func(in io.Reader) JSONDecoder
	}

	JSONEncoder interface {
		Encode(v any) error
	}

	JSONDecoder interface {
		More() bool
		RuneToken() (token rune, others any, err error)
		Decode(v any) error
	}
)

var DefaultConfig = StdConfig
