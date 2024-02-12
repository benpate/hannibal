package streams

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
)

// SetString sets a string property on the document
func (document *Document) SetString(name string, value ...string) bool {

	head := document.value.Head()

	if len(value) == 1 {
		document.value = head.Set(name, value[0])
		return true
	}

	document.value = head.Set(name, value)
	return true
}

// Append appends a value to a property on the document
func (document *Document) Append(name string, value any) bool {

	// RULE: If the value is empty, then NOOP
	if value == "" {
		return false
	}

	currentValue := convert.SliceOfAny(document.value.Head().Get(name))
	newValue := append(currentValue, value)

	document.value.Set(name, newValue)
	return true
}

// AppendString appends a string to a property on the document
func (document *Document) AppendString(name string, value string) bool {

	// RULE: If the value is empty, then NOOP
	if value == "" {
		return false
	}

	currentValue := convert.SliceOfAny(document.value.Head().Get(name))
	newValue := append(currentValue, value)

	document.value.Set(name, newValue)
	return true
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
