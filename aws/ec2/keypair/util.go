package keypair

import (
	"log"
	"strconv"
)

// StringValue gets the value for `key` as a string from the underlying map.
func StringValue(properties map[string]interface{}, key, defaultValue string) (value string) {
	value = defaultValue
	if val, ok := properties[key]; ok {
		if str, ok := val.(string); !ok {
			log.Printf("%s: unexpected type `%T`", key, val)
		} else {
			value = str
		}
	}
	return
}

// BoolValue gets the value for `key` as a string from the underlying map.
func BoolValue(properties map[string]interface{}, key string, defaultValue bool) (value bool) {
	value = defaultValue
	if val, ok := properties[key]; ok {
		switch t := val.(type) {
		case bool:
			value = t
		case string:
			value, _ = strconv.ParseBool(t)
		default:
			log.Printf("%s: unexpected type `%T`", key, val)
		}
	}
	return
}
