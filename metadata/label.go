package metadata

// Label is one moderation result applied to a document for the current viewer. A Label either
// hides the document (a block or a mute) or merely annotates it; IsHidden tells them apart.
// (Moderation systems often use "label" for the annotating kind alone; here it is the umbrella
// for every rule effect.) Value and Href are for DISPLAY ONLY -- never branch on their contents,
// or a reworded or localized string silently changes behavior.
type Label struct {
	Value    string // Human-readable text (e.g. "Muted by your rules", "Politics")
	Href     string // Optional link to more detail: the rule, a label definition, or the label's source
	IsHidden bool   // TRUE if this Label hides the document (a block or a mute); FALSE for an annotation
}
