# Hannibal

<img src="https://github.com/benpate/hannibal/raw/main/meta/logo.jpg">

Hannibal is an experimental ActivityPub library for Go. There are other packages/frameworks out there that are more complete and mature. So please check out [go-fed](https://github.com/go-fed) and [go-ap](https://github.com/go-ap) before trying Hannibal.

I'm writing this because the existing frameworks are very sophisticated, yet difficult for me to use.  They contain too much *magic* to be easily understood.  Hannibal fits the rest of the coding style of [Emissary](https://github.com/EmissarySocial/emissary) well enough to integrate nicely.

Like the ActivityPub spec itself, Hannibal is broken into several layers:

### ActivityPub client/server [pub]
https://www.w3.org/TR/activitypub/

This is not an ActivityPub framework, but a simple library that easily plugs into your existing app.  Add ActivityPub behaviors to your existing handlers, and send ActivityPub messages to 

### ActivityStream Vocabulary [vocab]
https://www.w3.org/TR/activitystreams-vocabulary/

### ActivityStreams data structures [stream]
https://www.w3.org/TR/activitystreams-core/

### JSON-LD reader/writer [jsonld]