package streams

import "encoding/json"

/******************************************
 * Custom JSON Marshalling / Unmarshalling
 ******************************************/

// MarshalJSON implements the json.Marshaller interface,
// and provides a custom marshalling into JSON --
// essentially just aiming the marshaller at the
// Document's value.
func (document Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(document.value)
}

// UnmarshalJSON implements the json.Unmarshaller interface,
// and provides a custom un-marshalling from JSON --
// essentially just aiming the unmashaller at the
// Document's value
func (document *Document) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &document.value)
}
