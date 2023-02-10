# Hannibal
ActivityPub library for Go

This is still experimental. There are other libraries out there that are far more mature. So please check out [go-fed](https://github.com/go-fed) and [go-ap](https://github.com/go-ap) before trying Hannibal.

I'm writing this because the existing frameworks are very sophisticated, yet difficult for me to use.  They contain too much *magic* to be easily understood.  Hannibal fits the rest of the coding style of [Emissary](https://github.com/EmissarySocial/emissary) well enough to integrate nicely.

Like the ActivityPub spec itself, Hannibal is broken into several layers:

### ActivityPub client/server [pub]
https://www.w3.org/TR/activitypub/

### ActivityStream Vocabulary [vocab]
https://www.w3.org/TR/activitystreams-vocabulary/

### ActivityStreams data structures [stream]
https://www.w3.org/TR/activitystreams-core/

### JSON-LD reader/writer [jsonld]