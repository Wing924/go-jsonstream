package jsonstream

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewJSONLinesEncoder_empty(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			sender := NewJSONLinesEncoderWithConfig(out, c.config)
			assert.NoError(t, sender.Done())
			assert.Equal(t, "", out.String())
		})
	}
}

func TestNewJSONLinesEncoder_small(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			sender := NewJSONLinesEncoderWithConfig(out, c.config)

			assert.NoError(t, sender.Send(smallSamples[0]))
			assert.Equal(t, `{"name":"Barbie"}
`, out.String())

			assert.NoError(t, sender.Send(smallSamples[1]))
			assert.Equal(t, `{"name":"Barbie"}
{"name":"Barbie"}
`, out.String())

			assert.NoError(t, sender.Send(smallSamples[2]))
			assert.Equal(t, `{"name":"Barbie"}
{"name":"Barbie"}
{"name":"Ken"}
`, out.String())

			assert.NoError(t, sender.Done())
			assert.Equal(t, `{"name":"Barbie"}
{"name":"Barbie"}
{"name":"Ken"}
`, out.String())
		})
	}
}

func TestNewJSONLinesDecoder_empty(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			var data Small
			{
				in := bytes.NewBufferString("")
				sender := NewJSONArrayDecoderWithConfig(in, c.config)
				assert.Equal(t, io.EOF, sender.Receive(&data))
			}
			{
				in := bytes.NewBufferString("\n")
				sender := NewJSONArrayDecoderWithConfig(in, c.config)
				assert.Equal(t, io.EOF, sender.Receive(&data))
			}
		})
	}
}

func TestNewJSONLinesDecoder_bad(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			var data Small
			in := bytes.NewBufferString("{")
			sender := NewJSONLinesDecoderWithConfig(in, c.config)
			assert.Error(t, sender.Receive(&data))
		})
	}
}

func TestNewJSONLinesDecoder_small(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			var data Small
			{
				in := bytes.NewBufferString(`{"name":"Barbie"}
{"name":"Barbie"}
{"name":"Ken"}
`)
				receiver := NewJSONLinesDecoderWithConfig(in, c.config)

				assert.NoError(t, receiver.Receive(&data))
				assert.Equal(t, Small{Name: "Barbie"}, data)

				assert.NoError(t, receiver.Receive(&data))
				assert.Equal(t, Small{Name: "Barbie"}, data)

				assert.NoError(t, receiver.Receive(&data))
				assert.Equal(t, Small{Name: "Ken"}, data)

				assert.Equal(t, io.EOF, receiver.Receive(&data))
			}
		})
	}
}

func BenchmarkNewJSONLinesDecoder(b *testing.B) {
	for _, c := range configs[1:] {
		b.Run(c.name, func(b *testing.B) {
			input := bytes.NewBuffer(nil)
			input.WriteString("[\n")
			for i := 0; i < b.N; i++ {
				if i > 0 {
					input.WriteString(",")
				}
				b, _ := json.Marshal(benchSample)
				input.Write(b)
			}
			input.WriteString("]\n")
			var data Small
			receiver := NewJSONArrayDecoderWithConfig(input, c.config)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = receiver.Receive(&data)
			}
		})
	}
}
