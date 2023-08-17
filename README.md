# go-jsonstream
JSON Stream (JSON Lines, JSON Array) Go Library

### Installation

```bash
go get github.com/Wing924/go-jsonstream
```

## How to use

### Encode

```go
// encode as JSON Lines stream
out := bytes.NewBuffer(nil)

sender := NewJSONLinesEncoder(out)
sender.Send("foo")
sender.Send(map[string]string{"bar": "baz"})
sender.Send(123)
sender.Done()

fmt.Println(out.String())
// Will output:
// "foo"
// {"bar":"baz"}
// 123
```

```go
// encode as JSON Array stream
out := bytes.NewBuffer(nil)

sender := NewJSONArrayEncoder(out)
sender.Send("foo")
sender.Send(map[string]string{"bar": "baz"})
sender.Send(123)
sender.Done()

fmt.Println(out.String())
// Will output:
// [
// "foo"
// ,{"bar":"baz"}
// ,123
// ]
```

### Decode

```go
// encode as JSON Lines stream
in := bytes.NewBufferString(`
{"name":"Barbie"}
{"name":"Barbie"}
{"name":"Ken"}
`)

var data struct{
    Name string `json:"name"`
}

receiver := jsonstream.NewJSONLinesDecoder(in)
for {
    err := receiver.Receive(&data)
    if err != nil {
        if err == io.EOF {
            break
        }
        panic(err)
    }
    fmt.Printf("%+v\n", data)
}

// will output
// {Name:Barbie}
// {Name:Barbie}
// {Name:Ken}
```

```go
// encode as JSON Array stream
in := bytes.NewBufferString(`[{"name":"Barbie"},{"name":"Barbie"},{"name":"Ken"}]`)

var data struct{
Name string `json:"name"`
}

receiver := jsonstream.NewJSONArrayDecoder(in)
for {
    err := receiver.Receive(&data)
    if err != nil {
        if err == io.EOF {
            break
        }
        panic(err)
    }
    fmt.Printf("%+v\n", data)
}

// will output
// {Name:Barbie}
// {Name:Barbie}
// {Name:Ken}
```

### Use [go-json](https://github.com/goccy/go-json) instead of `encoding/json`

The `go-json` is fast JSON encoder/decoder compatible with `encoding/json` for Go.

```go
// Change the global
jsonstream.DefaultConfig = jsonstream.GoJSONConfig

// now use go-json to encode/decode JSON
jsonstream.NewJSONArrayEncoder(out)
jsonstream.NewJSONArrayDecoder(in)
jsonstream.NewJSONLinesEncoder(out)
jsonstream.NewJSONLinesDecoder(in)
```

```go
// Change individual
jsonstream.NewJSONArrayEncoderWithConfig(out, jsonstream.GoJSONConfig)
jsonstream.NewJSONArrayDecoderWithConfig(in, jsonstream.GoJSONConfig)
jsonstream.NewJSONLinesEncoderWithConfig(out, jsonstream.GoJSONConfig)
jsonstream.NewJSONLinesDecoderWithConfig(in, jsonstream.GoJSONConfig)
```