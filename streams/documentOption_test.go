package streams

import (
	"testing"

	"github.com/benpate/hannibal/metadata"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWithNoStore confirms that the option marks the document, and that a document is unmarked
// without it.
func TestWithNoStore(t *testing.T) {

	assert.True(t, NewDocument(nil, WithNoStore()).Metadata.NoStore)
	assert.False(t, NewDocument(nil).Metadata.NoStore)
}

// TestWithNoStore_AddOptionsLeavesOriginal confirms that marking a copy through AddOptions leaves
// the original document unmarked.
func TestWithNoStore_AddOptionsLeavesOriginal(t *testing.T) {

	original := NewDocument(map[string]any{vocab.PropertyID: "https://example.com/note"})
	marked := original.AddOptions(WithNoStore())

	assert.True(t, marked.Metadata.NoStore)
	assert.False(t, original.Metadata.NoStore)
}

// TestWithLabels confirms that the option attaches the labels.
func TestWithLabels(t *testing.T) {

	labels := metadata.LabelSet{{Value: "Muted", IsHidden: true}}
	document := NewDocument(nil, WithLabels(labels))

	assert.Equal(t, labels, document.Metadata.Labels)
	assert.True(t, document.Metadata.IsRuleHidden())
}

// TestWithLabels_CallerCannotChangeThem confirms that changing the caller's labels after the option
// is built does not change what the option attaches.
func TestWithLabels_CallerCannotChangeThem(t *testing.T) {

	labels := metadata.LabelSet{{Value: "Muted", IsHidden: true}}
	option := WithLabels(labels)

	labels[0].Value = "changed"

	assert.Equal(t, "Muted", NewDocument(nil, option).Metadata.Labels[0].Value)
}

// TestWithLabels_DocumentsDoNotShare confirms that two documents given the same option each own
// their labels.
func TestWithLabels_DocumentsDoNotShare(t *testing.T) {

	option := WithLabels(metadata.LabelSet{{Value: "Muted", IsHidden: true}})

	first := NewDocument(nil, option)
	second := NewDocument(nil, option)

	first.Metadata.Labels[0].Value = "changed"

	assert.Equal(t, "Muted", second.Metadata.Labels[0].Value)
}

// TestWithLabels_Nil confirms that nil labels leave the document without labels.
func TestWithLabels_Nil(t *testing.T) {
	assert.Nil(t, NewDocument(nil, WithLabels(nil)).Metadata.Labels)
}

// TestWithRelation confirms that the option sets the relation type and the related URL together.
func TestWithRelation(t *testing.T) {

	document := NewDocument(nil, WithRelation(vocab.RelationTypeReply, "https://example.com/parent"))

	assert.Equal(t, vocab.RelationTypeReply, document.Metadata.RelationType)
	assert.Equal(t, "https://example.com/parent", document.Metadata.RelationHref)
	assert.True(t, document.Metadata.HasRelationship())
}

// TestWithDocumentCategory confirms that the option sets the category.
func TestWithDocumentCategory(t *testing.T) {

	document := NewDocument(nil, WithDocumentCategory(vocab.ObjectTypeNote))

	assert.Equal(t, vocab.ObjectTypeNote, document.Metadata.DocumentCategory)
}

// TestWithMetadata_ReplacesFieldOptions pins the order rule: WithMetadata replaces all metadata, so
// a field option applied before it is lost, and one applied after it is kept.
func TestWithMetadata_ReplacesFieldOptions(t *testing.T) {

	stored := metadata.Metadata{HashedID: "abc123"}

	lost := NewDocument(nil, WithNoStore(), WithMetadata(stored))
	require.Equal(t, "abc123", lost.Metadata.HashedID)
	assert.False(t, lost.Metadata.NoStore)

	kept := NewDocument(nil, WithMetadata(stored), WithNoStore())
	require.Equal(t, "abc123", kept.Metadata.HashedID)
	assert.True(t, kept.Metadata.NoStore)
}
