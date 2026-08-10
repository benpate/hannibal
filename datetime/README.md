## Hannibal / datetime

This package formats timestamps for ActivityStreams 2.0 date-time properties — `published`, `updated`, `startTime`, `endTime`, `deleted`, and anything else the vocabulary types as `xsd:dateTime`. Use it for every date-time value you write into an ActivityStreams document.

[AS2 Core section 2.3](https://www.w3.org/TR/activitystreams-core/#dates) requires these values to conform to the RFC3339 `date-time` production. Output is always normalized to UTC, so it ends in `Z` rather than a numeric offset.

### Usage

```go
datetime.Now()                        // "2026-08-10T17:04:05Z"
datetime.Format(someTime)             // from a time.Time
datetime.FromUnixSeconds(1767366245)  // from a Unix epoch, in seconds
datetime.FromUnixMilli(1767366245000) // from a Unix epoch, in milliseconds
```

### Empty values

Every function returns an **empty string** for a value it cannot render conformantly — a zero time, or a year outside `[1, 9999]`. This is not an error condition; it is how you say "no date."

Omit the property when you get one:

```go
message := mapof.Any{
    vocab.PropertyType: vocab.ObjectTypeNote,
}

if published := datetime.Format(note.PublishDate); published != "" {
    message[vocab.PropertyPublished] = published
}
```

An absent `published` property is well-defined in AS2. A `published` of `1970-01-01T00:00:00Z` is not — it is a claim that the object was published at the Unix epoch, and peers will sort your document accordingly.

### Writing only

This package serializes; it does not parse. Inbound timestamps from remote servers are handled by [rosetta](https://github.com/benpate/rosetta)'s `convert.TimeWithLocale`, which accepts RFC3339 alongside the HTTP and RFC822 formats that peers emit in practice. Strict on write, lenient on read.

This is also **not** the format for HTTP header dates. The `Date` header of a signed request uses IMF-fixdate, which is [`sigs`](../sigs)' business and is not exported.
