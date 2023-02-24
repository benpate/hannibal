package signatures

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {

	header := `keyId="http://localhost/@63810ba2721f7a33807f25c7/pub/key",algorithm="hs2019",headers="(request-target) host date",signature="Kcn3TRZLmwIAsEY3HGL5gP4BKU5uJ9Qhm2kBSxjLOyG2//soTSAuld1RbNhlm4Fz9twgt5GlkP1Ppeno1mP8VzhM4KhhXnr9a8Bd+Ti9gIv4EHeJeKRCURy6PTMqsbbGR37s3Nif/UQAjUWy2q63klINeB+L4W1g38vTxu1m9NiCXBdgMuyFHbXhxnFzpNstFXxfr+Mx9+rbmxUWtqMKj4Zh7dedqD0g/nDPTPbYFKNWOgYfBFwbyl3gRrWgavkyP7L/qeNYnkB98ACWdQwGGJQ1Fi4JQQxXxtAVbOKgBRFMLGn/jH8yYxzrW5OILznvJS7pzFgYsnzomVbaAdYkSA=="`
	result := ParseSignatureHeader(header)

	require.Equal(t, "http://localhost/@63810ba2721f7a33807f25c7/pub/key", result["keyId"])
	require.Equal(t, "hs2019", result["algorithm"])
	require.Equal(t, "(request-target) host date", result["headers"])
	require.Equal(t, "Kcn3TRZLmwIAsEY3HGL5gP4BKU5uJ9Qhm2kBSxjLOyG2//soTSAuld1RbNhlm4Fz9twgt5GlkP1Ppeno1mP8VzhM4KhhXnr9a8Bd+Ti9gIv4EHeJeKRCURy6PTMqsbbGR37s3Nif/UQAjUWy2q63klINeB+L4W1g38vTxu1m9NiCXBdgMuyFHbXhxnFzpNstFXxfr+Mx9+rbmxUWtqMKj4Zh7dedqD0g/nDPTPbYFKNWOgYfBFwbyl3gRrWgavkyP7L/qeNYnkB98ACWdQwGGJQ1Fi4JQQxXxtAVbOKgBRFMLGn/jH8yYxzrW5OILznvJS7pzFgYsnzomVbaAdYkSA==", result["signature"])
}
