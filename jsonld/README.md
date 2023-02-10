# JSON-LD

This is a simple library for reading and writing JSON-LD documents in Go.

JSON-LD is an extremely flexible data format, which makes it difficult to use in a strictly typed language like Go.  This package wraps values like `map[string]any` and `[]any` in a "Reader" interface, which lets us use them without too much pain.

It also retrieves up missing documents automatically via HTTP(s) if you try to access properties that are not present in the original document.

There are many other, better packages that handle JSON-LD in a far more rigorous way, such as https://github.com/kazarena/json-gold and https://github.com/go-ap/jsonld.  You should probably use them instead.  This library exists because I needed a small, specific sub-set of features for Hannibal, and it was easier to have control of this feature set internally instead of adding another external dependency.