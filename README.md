# Hannibal

<img src="https://github.com/benpate/hannibal/raw/main/meta/logo.jpg" style="width:100%; display:block; margin-bottom:20px;"  alt="Oil painting titled: Hannibal in the Alps, by R.B. Davis">

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/hannibal)
[![Version](https://img.shields.io/github/v/release/benpate/hannibal?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/hannibal/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/hannibal/go.yml?style=flat-square)](https://github.com/benpate/hannibal/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/hannibal?style=flat-square)](https://goreportcard.com/report/github.com/benpate/hannibal)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/hannibal.svg?style=flat-square)](https://codecov.io/gh/benpate/hannibal)

## Triumphant ActivityPub for Go
Hannibal is an experimental ActivityPub library for Go. It's goal is to be a robust, idiomatic, and thoroughly documented ActivityPub implementation fits into your application without any magic or drama.

## DO NOT USE

This project is a work-in-progress, and should NOT be used by ANYONE, for ANY PURPOSE, under ANY CIRCUMSTANCES.  It is WILL BE CHANGED UNDERNEATH YOU WITHOUT NOTICE OR HESITATION, and is expressly GUARANTEED to blow up your computer, send your cat into an infinite loop, and combine your hot and cold laundry into a single cycle.

There are other packages/frameworks out there that are more complete and mature. So please check out [go-fed](https://github.com/go-fed) and [go-ap](https://github.com/go-ap) before trying this.


## Packages
Like the ActivityPub spec itself, Hannibal is broken into several layers:

### pub - ActivityPub client/server
https://www.w3.org/TR/activitypub/

This is not an ActivityPub framework, but a simple library that easily plugs into your existing app.  Add ActivityPub behaviors to your existing handlers, and send ActivityPub messages to 

### vocab - ActivityStreams Vocabulary
https://www.w3.org/TR/activitystreams-vocabulary/

The `vocab` package includes the standard ActivityStream vocabulary, including names of actions, objects and properties used in ActivityPub. 

### streams - ActivityStreams data structures
https://www.w3.org/TR/activitystreams-core/

The `streams` package contains common data structures defined in the ActivityStreams spec, notably definitions for: `Document`, `Collection`, `OrderedCollection`, `CollectionPage`, and `OrderedCollectionPage`.  These are used by ActivityPub to send and receive multiple records in one HTTP request.

This package also includes a lightweight wrapper around generic data structures (like `map[string]any` and `[]any`) that makes it easy to access data structures within an ActivityStreams/JSON-LD document.

### sigs - HTTP Signatures and Digests
https://datatracker.ietf.org/doc/draft-ietf-httpbis-message-signatures

The `sigs` package creates and verifies HTTP signatures and Digests.


## Pull Requests Welcome

This library is a work in progress, and will benefit from your experience reports, use cases, and contributions.  If you have an idea for making this library better, send in a pull request.  We're all in this together! 🐘
