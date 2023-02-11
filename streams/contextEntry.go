package streams

import (
	"encoding/json"

	"github.com/benpate/rosetta/mapof"
)

// ContextEntry
// https://www.w3.org/TR/json-ld/#the-context
type ContextEntry struct {
	Vocabulary string       // The primary vocabulary represented by the context/document.
	Language   string       // The language
	Extensions mapof.String // a map of additional namespaces that are included in this context/document.
}

func NewContextEntry(vocabulary string) ContextEntry {
	return ContextEntry{
		Vocabulary: vocabulary,
		Language:   "und",
		Extensions: mapof.NewString(),
	}
}

func (entry *ContextEntry) WithLanguage(language string) *ContextEntry {
	entry.Language = language
	return entry
}

func (entry *ContextEntry) WithExtension(key string, value string) *ContextEntry {
	if len(entry.Extensions) == 0 {
		entry.Extensions = mapof.NewString()
	}

	entry.Extensions[key] = value
	return entry
}

func (entry ContextEntry) MarshalJSON() ([]byte, error) {

	// If this context only has a vocabulary, then
	// use the short-form "string only" syntax
	if entry.IsVocabularyOnly() {
		return json.Marshal(entry.Vocabulary)
	}

	// Otherwise, use the long-form syntax as a JSON object
	result := mapof.NewAny()
	result["@vocab"] = entry.Vocabulary

	if entry.IsLanguageDefined() {
		result["@language"] = entry.Language
	}

	if entry.HasExtensions() {
		for key, value := range entry.Extensions {
			result[key] = value
		}
	}

	return json.Marshal(result)
}

func (entry ContextEntry) IsVocabularyOnly() bool {
	if entry.IsLanguageDefined() {
		return false
	}

	if entry.HasExtensions() {
		return false
	}

	return true
}

func (entry ContextEntry) IsLanguageDefined() bool {
	if entry.Language == "" {
		return false
	}

	if entry.Language == "und" {
		return false
	}

	return true
}

func (entry ContextEntry) HasExtensions() bool {
	return len(entry.Extensions) > 0
}
