package streams

// OptionStripContext instructs the Document.Map() method to remove the "@context" property from its ouput.
const OptionStripContext = "STRIP_CONTEXT"

// OptionStripRecipients instructs the Document.Map() method to remove all recipient properties from its output.
// (To, BTo, CC, BCC)
const OptionStripRecipients = "STRIP_RECCIPIENTS"
