package jsonstream

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewJSONArrayEncoder_empty(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			sender := NewJSONArrayEncoderWithConfig(out, c.config)
			assert.NoError(t, sender.Done())
			assert.Equal(t, "[]\n", out.String())
		})
	}
}

func TestNewJSONArrayEncoder_small(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			sender := NewJSONArrayEncoderWithConfig(out, c.config)

			assert.NoError(t, sender.Send(smallSamples[0]))
			assert.Equal(t, `[
{"name":"Barbie"}
`, out.String())

			assert.NoError(t, sender.Send(smallSamples[1]))
			assert.Equal(t, `[
{"name":"Barbie"}
,{"name":"Barbie"}
`, out.String())

			assert.NoError(t, sender.Send(smallSamples[2]))
			assert.Equal(t, `[
{"name":"Barbie"}
,{"name":"Barbie"}
,{"name":"Ken"}
`, out.String())

			assert.NoError(t, sender.Done())
			assert.Equal(t, `[
{"name":"Barbie"}
,{"name":"Barbie"}
,{"name":"Ken"}
]
`, out.String())
		})
	}
}

func TestNewJSONArrayDecoder_empty(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			var data Small
			{
				in := bytes.NewBufferString("")
				sender := NewJSONArrayDecoderWithConfig(in, c.config)
				assert.Equal(t, io.EOF, sender.Receive(&data))
			}
			{
				in := bytes.NewBufferString("[]")
				sender := NewJSONArrayDecoderWithConfig(in, c.config)
				assert.Equal(t, io.EOF, sender.Receive(&data))
			}
		})
	}
}

func TestNewJSONArrayDecoder_bad(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name+"_[", func(t *testing.T) {
			var data Small
			{
				in := bytes.NewBufferString("[")
				sender := NewJSONArrayDecoderWithConfig(in, c.config)
				assert.ErrorIs(t, sender.Receive(&data), ErrBadJSONArray)
			}
		})
		t.Run(c.name+"_[{", func(t *testing.T) {
			var data Small
			{
				in := bytes.NewBufferString(`[{"name":"Ken"]`)
				sender := NewJSONArrayDecoderWithConfig(in, c.config)
				assert.Error(t, sender.Receive(&data))
			}
		})
	}
}

func TestNewJSONArrayDecoder_small(t *testing.T) {
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			var data Small
			{
				in := bytes.NewBufferString(`[{"name":"Barbie"},{"name":"Barbie"},{"name":"Ken"}]`)
				receiver := NewJSONArrayDecoderWithConfig(in, c.config)

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
