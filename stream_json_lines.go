package jsonstream

import (
	"io"
)

type (
	jsonLinesSender struct {
		writer  io.Writer
		encoder JSONEncoder
	}

	jsonLinesReceiver struct {
		reader  io.Reader
		decoder JSONDecoder
	}
)

func NewJSONLinesEncoder(out io.Writer) Sender {
	return NewJSONLinesEncoderWithConfig(out, DefaultConfig)
}

func NewJSONLinesEncoderWithConfig(out io.Writer, cfg *Config) Sender {
	return &jsonLinesSender{
		writer:  out,
		encoder: cfg.NewEncoder(out),
	}
}

func (s *jsonLinesSender) Send(data any) error {
	return s.encoder.Encode(data)
}

func (s *jsonLinesSender) Done() error {
	return nil
}

func NewJSONLinesDecoder(in io.Reader) Receiver {
	return NewJSONLinesDecoderWithConfig(in, DefaultConfig)
}

func NewJSONLinesDecoderWithConfig(in io.Reader, cfg *Config) Receiver {
	return &jsonLinesReceiver{
		reader:  in,
		decoder: cfg.NewDecoder(in),
	}
}

func (r *jsonLinesReceiver) Receive(data any) error {
	return r.decoder.Decode(data)
}
