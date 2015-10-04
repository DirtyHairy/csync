package factory

import (
	"encoding/json"
	"testing"

	"github.com/DirtyHairy/csync/lib/storage/config"
	config_local "github.com/DirtyHairy/csync/lib/storage/local/config"
	"github.com/DirtyHairy/csync/testutils"
	"github.com/blang/semver"
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
