package jsonstream

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrBadJSONArray = errors.New("bad JSON Array")
)

type (
	jsonArraySender struct {
		writer    io.Writer
		encoder   JSONEncoder
		committed bool
	}

	jsonArrayReceiver struct {
		reader    io.Reader
		decoder   JSONDecoder
		committed bool
	}
)

var (
	_ Sender   = (*jsonArraySender)(nil)
	_ Receiver = (*jsonArrayReceiver)(nil)
)

func NewJSONArrayEncoder(out io.Writer, cfg *Config) Sender {
	return NewJSONArrayEncoderWithConfig(out, DefaultConfig)
}

func NewJSONArrayEncoderWithConfig(out io.Writer, cfg *Config) Sender {
	if cfg == nil {
		cfg = DefaultConfig
	}
	return &jsonArraySender{
		writer:  out,
		encoder: cfg.NewEncoder(out),
	}
}

func (s *jsonArraySender) Send(data any) error {
	if !s.committed {
		s.committed = true
		_, err := s.writer.Write([]byte("[\n"))
		if err != nil {
			return err
		}
	} else if _, err := s.writer.Write([]byte(",")); err != nil {
		return err
	}

	return s.encoder.Encode(data)
}

func (s *jsonArraySender) Done() error {
	if !s.committed {
		_, err := s.writer.Write([]byte("[]\n"))
		return err
	}
	_, err := s.writer.Write([]byte("]\n"))
	if err != nil {
		return err
	}
	s.committed = false
	return nil
}

func NewJSONArrayDecoder(in io.Reader) Receiver {
	return NewJSONArrayDecoderWithConfig(in, DefaultConfig)
}

func NewJSONArrayDecoderWithConfig(in io.Reader, cfg *Config) Receiver {
	if cfg == nil {
		cfg = DefaultConfig
	}
	return &jsonArrayReceiver{
		reader:  in,
		decoder: cfg.NewDecoder(in),
	}
}

func (r *jsonArrayReceiver) Receive(data any) error {
	if !r.committed {
		r.committed = true
		token, others, err := r.decoder.RuneToken()
		if err != nil {
			return err
		}
		if others != nil {
			return fmt.Errorf("%w: expect '[' but got %T(%v)", ErrBadJSONArray, others, others)
		}
		if token != '[' {
			return fmt.Errorf("%w: expect '[' but got %q", ErrBadJSONArray, token)
		}
	}
	if !r.decoder.More() {
		token, others, err := r.decoder.RuneToken()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("%w: expect ']' but got EOF", ErrBadJSONArray)
			}
			return err
		}
		if others != nil {
			return fmt.Errorf("%w: expect '[' but got %T(%v)", ErrBadJSONArray, others, others)
		}
		if token != ']' {
			return fmt.Errorf("%w: expect '[' but got %q", ErrBadJSONArray, token)
		}
		return io.EOF
	}
	if err := r.decoder.Decode(&data); err != nil {
		return err
	}
	return nil
}
