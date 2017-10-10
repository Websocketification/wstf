package wstf

import "encoding/json"

// Stringify object.
func Stringify(m interface{}) string {
	str, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(str)
}
