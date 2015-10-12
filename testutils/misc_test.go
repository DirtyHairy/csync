package testutils

import (
	"testing"
)

func testJsonEqual(json1, json2 string, t *testing.T) {
	w := Wrap(t)

	isEqual, err := JsonEqual([]byte(json1), []byte(json2))

	w.BailIfErrorf(err, "JSON failed to deserialize: %v", err)

	if !isEqual {
		t.Fatalf("JSON failed to compare as equal: %s vs. %s", json1, json2)
	}
}

func testJsonNotEqual(json1, json2 string, t *testing.T) {
	w := Wrap(t)

	isEqual, err := JsonEqual([]byte(json1), []byte(json2))

	w.BailIfErrorf(err, "JSON failed to deserialize: %v", err)

	if isEqual {
		t.Fatalf("JSON compared as equal: %s vs. %s", json1, json2)
	}
}

func TestJsonEqual(t *testing.T) {
	testJsonEqual(`1`, `1  `, t)
	testJsonNotEqual(`1`, `2`, t)

	testJsonEqual(`true`, ` true`, t)
	testJsonNotEqual(`false`, `true`, t)

	testJsonEqual(`"huppe"`, `"huppe" `, t)
	testJsonNotEqual(`"huppe"`, `"hoppe"`, t)

	testJsonEqual(`
        {
            "foo": 1,
            "bar": true
        }
    `, `
        {
            "bar": true,
            "foo": 1
        }
    `, t)
	testJsonNotEqual(`
        {
            "foo": 1,
            "bar": true
        }
    `, `
        {
            "foo": true,
            "bar": 1
        }
    `, t)

	testJsonEqual(`
        {
            "foo": [1, 2, 3]
        }
    `, `
        {
            "foo": [1,
                2,
                    3]
        }
    `, t)
	testJsonNotEqual(`
        {
            "foo": [1, 3, 2]
        }
    `, `
        {
            "foo": [1, 2, 3]
        }
    `, t)
}
