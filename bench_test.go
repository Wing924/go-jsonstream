package jsonstream

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type Format int

const (
	FormatJSONLines Format = iota
	FormatJSONArray
)

var formatNames = map[Format]string{
	FormatJSONLines: "JSON Lines",
	FormatJSONArray: "JSON Array",
}

func (f Format) String() string {
	return formatNames[f]
}

func NewEncoder(out io.Writer, format Format, cfg *Config) Sender {
	switch format {
	case FormatJSONLines:
		return NewJSONLinesEncoderWithConfig(out, cfg)
	case FormatJSONArray:
		return NewJSONArrayEncoderWithConfig(out, cfg)
	default:
		panic(fmt.Errorf("unknown format: %v", format))
	}
}

func NewDecoder(in io.Reader, format Format, cfg *Config) Receiver {
	switch format {
	case FormatJSONLines:
		return NewJSONLinesDecoderWithConfig(in, cfg)
	case FormatJSONArray:
		return NewJSONArrayDecoderWithConfig(in, cfg)
	default:
		panic(fmt.Errorf("unknown format: %v", format))
	}
}

func BenchmarkEncoder(b *testing.B) {
	for _, c := range configs[1:] {
		for _, f := range []Format{FormatJSONLines, FormatJSONArray} {
			b.Run(c.name+"_"+f.String(), func(b *testing.B) {
				sender := NewEncoder(io.Discard, f, c.config)
				for i := 0; i < b.N; i++ {
					_ = sender.Send(benchSample)
				}
				_ = sender.Done()
			})
		}
	}
}

func BenchmarkDecoder(b *testing.B) {
	for _, c := range configs[1:] {
		for _, f := range []Format{FormatJSONLines, FormatJSONArray} {
			b.Run(c.name+"_"+f.String(), func(b *testing.B) {
				input := bytes.NewBuffer(nil)
				sender := NewEncoder(input, f, c.config)
				for i := 0; i < b.N; i++ {
					_ = sender.Send(benchSample)
				}
				_ = sender.Done()

				var data Small
				receiver := NewDecoder(input, f, c.config)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = receiver.Receive(&data)
				}

			})
		}
	}
}
