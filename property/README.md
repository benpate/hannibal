## Hannibal / unit

This package wraps a number of common data values with conversions used in ActivityStreams.  This is important because JSON-LD allows any piece of data to be represented in multiple formats.  This wrapper allows safe access to indidividual values even when there are nil values  slices, or maps of values present.

### Usage
```go
value := property.NewValue("http://foo.com")

foo.Raw() // returns "foo"

// Traverse Arrays
foo.Len() // returns 1
foo.Head().Raw() // returns "http://foo.com"
foo.Tail().Raw() // returns a nil value

// Represent maps
foo.Map() // returns a map with "id"="http://foo.com"
foo.Set("name", "Foo") // converts foo to a map and sets a new property
foo.Get("name").Raw() // returns "Foo"
```

## Interfaces
Everything implements the `property.Value` interface, which provides several low-level manipulations for reading, writing, and transforming values.  You can implement this interface in other packages to use other custom types with the rest of the hannibal library.
