package streams

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
)

// SetString sets a string property on the document
func (document *Document) SetString(name string, value ...string) bool {

	document.forceValueToMap()

	if asMap, ok := document.value.(map[string]any); ok {

		switch len(value) {
		case 0:
			delete(asMap, name)
		case 1:
			asMap[name] = value[0]
		default:
			asMap[name] = value
		}
		return true
	}

	return false
}

// AppendString appends a string to a property on the document
func (document *Document) AppendString(name string, value string) bool {

	// RULE: If the value is empty, then NOOP
	if value == "" {
		return false
	}

	document.forceValueToMap()

	if asMap, ok := document.value.(map[string]any); ok {
		oldValue := convert.SliceOfString(asMap[name])
		newValue := append(oldValue, value)
		asMap[name] = newValue
		return true
	}

	return false
}

// forceValueToMap forces the value into a Map format
func (document *Document) forceValueToMap() {

	switch typed := document.value.(type) {
	case string:
		document.value = map[string]any{
			vocab.PropertyID: typed,
		}

	case []any:
		if len(typed) == 0 {
			document.value = map[string]any{}
		} else {
			document.value = map[string]any{
				vocab.PropertyID: typed[0],
			}
		}
	}
}

/******************************************
 * Set/Append Helpers
 ******************************************/

// SetTo sets the "to" property of the document
func (document *Document) SetTo(value ...string) bool {
	return document.SetString(vocab.PropertyTo, value...)
}

// SetBto sets the "bto" property of the document
func (document *Document) SetBTo(value ...string) bool {
	return document.SetString(vocab.PropertyBTo, value...)
}

// SetCC sets the "cc" property of the document
func (document *Document) SetCC(value ...string) bool {
	return document.SetString(vocab.PropertyCC, value...)
}

// SetBCC sets the "bcc" property of the document
func (document *Document) SetBCC(value ...string) bool {
	return document.SetString(vocab.PropertyBCC, value...)
}

// AppendTo appends a value to the "to" property of the document
func (document *Document) AppendTo(value string) bool {
	return document.AppendString(vocab.PropertyTo, value)
}

// AppendBTo appends a value to the "bto" property of the document
func (document *Document) AppendBTo(value string) bool {
	return document.AppendString(vocab.PropertyBTo, value)
}

// AppendCC appends a value to the "cc" property of the document
func (document *Document) AppendCC(value string) bool {
	return document.AppendString(vocab.PropertyCC, value)
}

// AppendBCC appends a value to the "bcc" property of the document
func (document *Document) AppendBCC(value string) bool {
	return document.AppendString(vocab.PropertyBCC, value)
}
