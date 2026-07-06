package sender

// Option is a functional option that configures a Sender at construction time.
type Option func(*Sender)

// AllowPrivateIPs returns an Option that permits (or forbids) outbound delivery to
// non-public (private/loopback) addresses. Leave this unset in production so that
// remote's SSRF guard stays active; enable it only when delivering to a local peer.
func AllowPrivateIPs(value bool) Option {
	return func(sender *Sender) {
		sender.allowPrivateIPs = value
	}
}
