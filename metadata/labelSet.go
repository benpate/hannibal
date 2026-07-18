package metadata

import "slices"

// LabelSet is every effect the current viewer's rules apply to one document. It is computed
// per-request by a rules client, attached to a document's Metadata at load time, and never
// persisted or serialized. The zero value (a nil set) means "no rules matched".
//
// By convention hidden entries sort before annotations, so Reason() can name the headline verdict
// -- but correctness always flows through these methods, never through element position.
type LabelSet []Label

// IsHidden returns TRUE if any matched rule hides this document (a block or a mute).
func (set LabelSet) IsHidden() bool {
	return slices.ContainsFunc(set, func(label Label) bool {
		return label.IsHidden
	})
}

// HasAnnotations returns TRUE if any matched rule annotates this document without hiding it.
func (set LabelSet) HasAnnotations() bool {
	return slices.ContainsFunc(set, func(label Label) bool {
		return !label.IsHidden
	})
}

// Reason returns the display text of the first hiding Label -- the headline for a placeholder --
// or "" when nothing hides this document.
func (set LabelSet) Reason() string {

	for _, label := range set {
		if label.IsHidden {
			return label.Value
		}
	}

	return ""
}

// Hidden returns just the Labels that hide this document.
func (set LabelSet) Hidden() LabelSet {
	return set.filter(true)
}

// Annotations returns just the Labels that annotate this document without hiding it.
func (set LabelSet) Annotations() LabelSet {
	return set.filter(false)
}

// filter returns the subset of Labels whose IsHidden matches the argument.
func (set LabelSet) filter(hidden bool) LabelSet {

	result := make(LabelSet, 0, len(set))

	for _, label := range set {
		if label.IsHidden == hidden {
			result = append(result, label)
		}
	}

	return result
}

// Clone returns a copy of this set that shares no backing array with the original.
func (set LabelSet) Clone() LabelSet {
	return slices.Clone(set)
}
