package local

import (
	"encoding/json"
	"testing"

	"github.com/DirtyHairy/csync/lib/storage/local/config"
	"github.com/DirtyHairy/csync/testutils"
	"github.com/blang/semver"
)

func TestUnmarshal(t *testing.T) {
	w := testutils.Wrap(t)

	var err error
	version := semver.MustParse("0.0.1")

	jsonData := []byte(`{
        "version": "0.0.1",
        "path": "./test_artifacts"
    }`)

	cfg := config.NewConfig(version)

	err = json.Unmarshal(jsonData, &cfg)
	w.BailIfErrorf(err, "deserializazion failed: %v", err)

	provider, err := FromConfig(cfg)
	w.BailIfErrorf(err, "instantiation failed: %v", err)

	root, err := provider.Root()
	w.BailIfError(err)

	_, err = checkDirectoryContents(root, expectedRootEntries())

	w.BailIfErrorf(err, "provider path differs")
}

func TestMarshalUnmarshal(t *testing.T) {
	w := testutils.Wrap(t)

	var err error
	version := semver.MustParse("0.0.1")

	provider := getFSInstace()
	cfg := provider.Marshal()

	jsonData, err := json.Marshal(cfg)
	w.BailIfErrorf(err, "serialization failed: %v", err)

	unmarshalCfg := config.NewConfig(version)

	err = json.Unmarshal(jsonData, &unmarshalCfg)
	w.BailIfErrorf(err, "deserialization failed: %v", err)

	unmarshalledProvider, err := FromConfig(unmarshalCfg)
	w.BailIfErrorf(err, "instatiation failed: %v", err)

	root, err := unmarshalledProvider.Root()
	w.BailIfError(err)

	_, err = checkDirectoryContents(root, expectedRootEntries())

	w.BailIfErrorf(err, "provider path differs")
}
