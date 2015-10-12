package testutils

import (
	"encoding/json"
	"reflect"
)

func JsonEqual(json1, json2 []byte) (isEqual bool, err error) {
	var deserialized1, deserialized2 interface{}

	err = json.Unmarshal(json1, &deserialized1)
	if err != nil {
		return
	}

	err = json.Unmarshal(json2, &deserialized2)
	if err != nil {
		return
	}

	isEqual = reflect.DeepEqual(deserialized1, deserialized2)
	return
}
