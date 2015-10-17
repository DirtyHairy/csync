package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/DirtyHairy/csync/lib/repository"
	"github.com/DirtyHairy/csync/lib/repository/registry/config"
	"github.com/DirtyHairy/csync/lib/storage/local"
	"github.com/DirtyHairy/csync/testutils"
	"github.com/blang/semver"
)

func TestUnmarshalWithLocalFS(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	version := semver.MustParse("0.0.1")

	jsonData := []byte(`{
        "foo": {
            "plain": true,
            "storage": {
                "type": "local",
                "path": "."
            }
        },
        "bar": {
            "plain": true,
            "storage": {
                "type": "local",
                "path": ".."
            }
        }
    }`)

	cfg := config.NewConfig(version)

	err = json.Unmarshal(jsonData, &cfg)
	w.BailIfErrorf(err, "deserialization failed: %v", err)

	registry, err := CreateFromConfig(cfg)
	w.BailIfErrorf(err, "unmarshalling failed: %v", err)

	if registry.Size() != 2 {
		t.Fatalf("registry has the wrong size; expected 2, got %d", registry.Size())
	}

	if registry.Get("foo") == nil {
		t.Fatalf("foo: no such key")
	}

	if registry.Get("bar") == nil {
		t.Fatalf("bar: no such key")
	}
}

func TestMarshalWithLocalFS(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	cwd, err := os.Getwd()
	w.BailIfError(err)

	referenceJson := []byte(fmt.Sprintf(`{
	    "foo": {
            "plain": true,
            "storage": {
                "type": "local",
                "path": "%s"
            }
        },
        "bar": {
            "plain": true,
            "storage": {
                "type": "local",
                "path": "%s"
            }
        }
	}`, cwd, cwd))

	registry := NewRegistry()

	storage1, err := local.NewLocalFS(".")
	w.BailIfErrorf(err, "unable to init storage: %v", err)

	storage2, err := local.NewLocalFS(".")
	w.BailIfErrorf(err, "unable to init storage: %v", err)

	repo1 := repository.NewRepository(true, storage1)
	repo2 := repository.NewRepository(true, storage2)

	registry.Set("foo", repo1)
	registry.Set("bar", repo2)

	cfg := registry.Marshal()
	jsonData, err := json.Marshal(cfg)
	w.BailIfErrorf(err, "serialization failed: %v", err)

	match, err := testutils.JsonEqual(jsonData, referenceJson)
	w.BailIfErrorf(err, "error parsing JSON: %v", err)

	if !match {
		t.Fatalf("JSON did not match:\n%s\nvs.\n%s", string(jsonData), string(referenceJson))
	}
}
