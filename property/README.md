## Hannibal / property

This package wraps common data values with the conversions used throughout ActivityStreams. JSON-LD allows the same piece of data to be represented in several shapes — a bare value, a single-element array, or a map keyed by `id` — and this wrapper gives you safe access to the underlying value no matter which shape it arrived in, even in the presence of `nil` values, slices, or nested maps.

### Usage

```go
value := property.NewValue("http://foo.com")

value.Raw() // returns "http://foo.com"

// Treat any value as an array
value.Len()         // returns 1
value.Head().Raw()  // returns "http://foo.com"
value.Tail().IsNil() // returns true (nothing after the first element)

// Treat any value as a map (JSON-LD represents a bare URL as {"id": "..."})
value.Map()                  // returns map[string]any{"id": "http://foo.com"}

// Set returns a new value with the property applied
named := value.Set("name", "Foo")
named.Get("name").Raw()      // returns "Foo"
```

## Interfaces

Everything implements the `property.Value` interface, which provides the low-level operations for reading, writing, and transforming values: `Get`, `Set`, `Head`, `Tail`, `Len`, `IsNil`, `Map`, and `Raw`. You can implement this interface in your own package to use custom types with the rest of the Hannibal library.
