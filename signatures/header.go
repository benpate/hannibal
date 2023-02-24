package signatures

import (
	"strings"

	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/mapof"
)

func ParseSignatureHeader(value string) mapof.String {

	result := mapof.NewString()

	item := ""
	itemList := list.ByComma(value)

	for !itemList.IsEmpty() {
		item, itemList = itemList.Split()
		name, value := list.Split(item, '=')
		value = strings.Trim(value, "\" ")

		result[name] = value
	}

	return result
}
