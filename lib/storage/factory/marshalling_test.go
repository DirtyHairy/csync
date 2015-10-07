package factory

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/DirtyHairy/csync/lib/storage/config"
	"github.com/DirtyHairy/csync/lib/storage/local"
	"github.com/DirtyHairy/csync/testutils"
	"github.com/blang/semver"

	config_local "github.com/DirtyHairy/csync/lib/storage/local/config"
)

func TestUnmarshalLocal(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	version := semver.MustParse("0.0.1")
	jsonData := `{
	           "type": "local",
	           "path": "."
	       }`

	cfg := config.NewConfig(version)
	err = json.Unmarshal([]byte(jsonData), &cfg)

	w.BailIfErrorf(err, "unmarshalling failed: %v", err)

	_, err = FromConfig(cfg)

	w.BailIfErrorf(err, "creating the provider failed: %v", err)
}

func TestUnmarshalLocalWithInvalidPath(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	version := semver.MustParse("0.0.1")
	jsonData := `{
	           "type": "local",
	           "path": "./INVALID_PATH"
	       }`

	cfg := config.NewConfig(version)
	err = json.Unmarshal([]byte(jsonData), &cfg)

	w.BailIfErrorf(err, "unmarshalling failed: %v", err)

	_, err = FromConfig(cfg)

	w.ShouldFailf(err, "unmarshalling should fail due to the invalid path")
}

func TestUnmarshalLocalWithInvalidType(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	version := semver.MustParse("0.0.1")
	jsonData := `{
	           "type": "localee",
	           "path": "."
	       }`

	cfg := config.NewConfig(version)
	err = json.Unmarshal([]byte(jsonData), &cfg)

	w.ShouldFailf(err, "unmarshalling should fail due to invalid type")
}

func TestUnmarshalLocalWithInvalidJSON(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	version := semver.MustParse("0.0.1")
	jsonData := `{
	           "type": "local",
	           "path": "."
	       `

	cfg := config.NewConfig(version)
	err = json.Unmarshal([]byte(jsonData), &cfg)

	w.ShouldFailf(err, "unmarshalling should fail due to invalid JSON")
}

func TestUnmarshalLocalFromSpecificConfig(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	version := semver.MustParse("0.0.1")
	jsonData := `{
	           "type": "local",
	           "path": "."
	        }`

	cfg := config_local.NewConfig(version)
	err = json.Unmarshal([]byte(jsonData), &cfg)

	w.BailIfErrorf(err, "unmarshalling failed: %v", err)

	_, err = FromConfig(cfg)

	w.BailIfErrorf(err, "creating the provider failed: %v", err)

}

func TestMarshalUnmarshal(t *testing.T) {
	w := testutils.Wrap(t)

	var err error
	version := semver.MustParse("0.0.1")

	provider, err := local.NewLocalFS("./test_artifacts")

	w.BailIfError(err)

	cfg := provider.Marshal()
	jsonData, err := json.Marshal(cfg)

	w.BailIfErrorf(err, "serialization failed: %v", err)

	if !strings.Contains(string(jsonData), "test_artifacts") {
		t.Fatalf("JSON does not contain the proper path: %s", string(jsonData))
	}

	unmarshalledConfig := config.NewConfig(version)
	err = json.Unmarshal(jsonData, &unmarshalledConfig)

	w.BailIfErrorf(err, "deserialization failed: %v", err)

	_, err = FromConfig(unmarshalledConfig)

	w.BailIfErrorf(err, "reinstation of fs provider failed: %v", err)
}
