package metadata

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// TestMetadata_IsRuleHidden confirms the Metadata delegate reflects its LabelSet, and that a
// zero-value Metadata reads as visible.
func TestMetadata_IsRuleHidden(t *testing.T) {

	assert.False(t, New().IsRuleHidden())

	hidden := Metadata{Labels: LabelSet{{Value: "Blocked", IsHidden: true}}}
	assert.True(t, hidden.IsRuleHidden())

	// An annotation-only set never hides.
	labeled := Metadata{Labels: LabelSet{{Value: "Politics", IsHidden: false}}}
	assert.False(t, labeled.IsRuleHidden())
}

// TestMetadata_Clone confirms a clone owns its own Labels slice.
func TestMetadata_Clone(t *testing.T) {

	original := Metadata{Labels: LabelSet{{Value: "Muted", IsHidden: true}}}
	clone := original.Clone()
	assert.Equal(t, original, clone)

	// Writing through the clone must not reach the original.
	clone.Labels[0].Value = "changed"
	assert.Equal(t, "Muted", original.Labels[0].Value)
}

// TestMetadata_Categories confirms the IsActor/IsObject/IsCollection predicates
// each match only their own DocumentCategory.
func TestMetadata_Categories(t *testing.T) {

	actor := Metadata{DocumentCategory: vocab.DocumentCategoryActor}
	assert.True(t, actor.IsActor())
	assert.False(t, actor.IsObject())
	assert.False(t, actor.IsCollection())

	object := Metadata{DocumentCategory: vocab.DocumentCategoryObject}
	assert.False(t, object.IsActor())
	assert.True(t, object.IsObject())
	assert.False(t, object.IsCollection())

	collection := Metadata{DocumentCategory: vocab.DocumentCategoryCollection}
	assert.False(t, collection.IsActor())
	assert.False(t, collection.IsObject())
	assert.True(t, collection.IsCollection())

	// A zero-value Metadata matches none of the categories.
	empty := New()
	assert.False(t, empty.IsActor())
	assert.False(t, empty.IsObject())
	assert.False(t, empty.IsCollection())
}

// TestMetadata_HasCounts confirms each Has* predicate is true only when its
// count is positive.
func TestMetadata_HasCounts(t *testing.T) {

	assert.True(t, Metadata{Replies: 1}.HasReplies())
	assert.False(t, Metadata{Replies: 0}.HasReplies())

	assert.True(t, Metadata{Announces: 1}.HasAnnounces())
	assert.False(t, Metadata{Announces: 0}.HasAnnounces())

	assert.True(t, Metadata{Likes: 1}.HasLikes())
	assert.False(t, Metadata{Likes: 0}.HasLikes())

	assert.True(t, Metadata{Dislikes: 1}.HasDislikes())
	assert.False(t, Metadata{Dislikes: 0}.HasDislikes())
}

// TestMetadata_HasRelationship confirms a relationship requires BOTH a type and
// an href.
func TestMetadata_HasRelationship(t *testing.T) {

	assert.True(t, Metadata{
		RelationType: vocab.RelationTypeReply,
		RelationHref: "https://example.com/1",
	}.HasRelationship())

	// Missing either half is not a relationship.
	assert.False(t, Metadata{RelationType: vocab.RelationTypeReply}.HasRelationship())
	assert.False(t, Metadata{RelationHref: "https://example.com/1"}.HasRelationship())
	assert.False(t, Metadata{}.HasRelationship())
}

// TestMetadata_SetRelationCount confirms the setter updates the matching counter
// and reports whether the value actually changed.
func TestMetadata_SetRelationCount(t *testing.T) {

	check := func(name string, relationType string, read func(Metadata) int64) {
		t.Run(name, func(t *testing.T) {
			metadata := Metadata{}

			// First write changes the value -> returns true.
			assert.True(t, metadata.SetRelationCount(relationType, 5))
			assert.Equal(t, int64(5), read(metadata))

			// Re-writing the same value is a no-op -> returns false.
			assert.False(t, metadata.SetRelationCount(relationType, 5))
		})
	}

	check("Reply", vocab.RelationTypeReply, func(m Metadata) int64 { return m.Replies })
	check("Announce", vocab.RelationTypeAnnounce, func(m Metadata) int64 { return m.Announces })
	check("Like", vocab.RelationTypeLike, func(m Metadata) int64 { return m.Likes })
	check("Dislike", vocab.RelationTypeDislike, func(m Metadata) int64 { return m.Dislikes })
}

// TestMetadata_SetRelationCount_Unknown confirms an unrecognized relation type
// changes nothing and reports false.
func TestMetadata_SetRelationCount_Unknown(t *testing.T) {

	metadata := Metadata{}
	assert.False(t, metadata.SetRelationCount("NotARealRelation", 99))
	assert.Equal(t, Metadata{}, metadata)
}
