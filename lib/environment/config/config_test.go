package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/DirtyHairy/csync/testutils"
	"github.com/blang/semver"
)

type mockLocator string

func (l mockLocator) Locate() (string, error) {
	return string(l), nil
}

func (l mockLocator) Destroy() {
	os.RemoveAll(string(l))
}

func newMockLocator() (l mockLocator, err error) {
	directory, err := ioutil.TempDir("", "csynt_test")

	if err != nil {
		return
	}

	l = mockLocator(directory)
	return
}

func TestMarshall(t *testing.T) {
	w := testutils.Wrap(t)
	version := semver.MustParse("1.2.3")

	var err error

	config := NewConfig(version)
	jsonData, err := marshalConfig(config)

	w.BailIfError(err)

	unmarshalledConfig, err := unmarshalConfig(jsonData)

	w.BailIfError(err)

	if actualVersion := unmarshalledConfig.CsyncVersion(); actualVersion.NE(version) {
		t.Fatalf("unmarshalled wrong version: expected %s, got %s", version, actualVersion)
	}
}

func TestUnmarshallEmpty(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	_, err = unmarshalConfig([]byte(`{}`))

	w.ShouldFailf(err, "unmarshalling should fail if version is missing")
}

func TestSave(t *testing.T) {
	w := testutils.Wrap(t)
	version := semver.MustParse("1.2.3")

	var err error

	locator, err := newMockLocator()
	w.BailIfError(err)
	defer locator.Destroy()

	manager := NewManager(locator)
	err = manager.Save(NewConfig(version))

	w.BailIfErrorf(err, "saving the config failed: %v", err)

	config, err := manager.Load()

	w.BailIfErrorf(err, "loading the config failed: %v", err)

	if actualVersion := config.CsyncVersion(); actualVersion.NE(version) {
		t.Fatalf("loaded wrong version: expected %s, got %s", version, actualVersion)
	}
}
