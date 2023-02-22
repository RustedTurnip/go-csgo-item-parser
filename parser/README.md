# parser

This package implements a parser specifically for the file formats present in the CSGO `items_game.txt`
and `csgo_<language>.txt` files. These use the Valve Data Format (VDF) and this parser is effectively
a VDF parser (*however only tested on the two aforementioned files*).

## Usage

Currently, to use this parser, the data must be located within a file and the location of that file can
be passed into the `Parse` function where it will be read and converted into a map of type
`map[string]interface{}`.

Example:

```go
result, err := parser.Parse("/path/to/file.txt")
```

The underlying data can be either a nested `map[string]interface{}` if the key holds a subsection of data,
or otherwise a `string` which represents the data.

e.g.

```vdf
"foo"
{
    "bar"
    {
        "foobar"    "one"
    }
}
```

Would be translated into:
```go
result := map[string]interface{}{
    "foo": map[string]interface {
        "bar": map[string]interface {
            "foobar": "one",			
        }
    }
}
```