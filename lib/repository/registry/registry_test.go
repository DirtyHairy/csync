package registry

import (
	"reflect"
	"sort"
	"testing"

	"github.com/DirtyHairy/csync/lib/repository"
	"github.com/DirtyHairy/csync/lib/storage/local"
	"github.com/DirtyHairy/csync/testutils"
)

func newLocalFSRepository(path string) (repo repository.Repository, err error) {
	storage, err := local.NewLocalFS(path)

	if err != nil {
		return
	}

	repo = repository.NewRepository(true, storage)

	return
}

func TestBasicUsage(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	repo, err := newLocalFSRepository(".")
	w.BailIfError(err)

	registry := NewRegistry()

	if registry.Size() != 0 {
		t.Fatalf("registry should start empty, got size: %d", registry.Size())
	}

	if registry.Get("foo") != nil {
		t.Fatalf("Getting a nonexisting key should return nil, got %v", registry.Get("foo"))
	}

	registry.Set("foo", repo)

	if registry.Size() != 1 {
		t.Fatalf("registry should grow to size 1 after set, got %d", registry.Size())
	}

	if registry.Get("foo") != repo {
		t.Fatalf("failed to retrieve repo")
	}

	registry.Unset("foo")
	if registry.Size() != 0 {
		t.Fatalf("clearing the key failed")
	}
}

func TestGetKeys(t *testing.T) {
	w := testutils.Wrap(t)

	var err error

	repo1, err := newLocalFSRepository(".")
	w.BailIfError(err)

	repo2, err := newLocalFSRepository("..")
	w.BailIfError(err)

	registry := NewRegistry()

	registry.Set("foo", repo1)
	registry.Set("bar", repo2)

	keys := registry.Keys()
	sort.Strings(keys)

	if !reflect.DeepEqual(keys, []string{"bar", "foo"}) {
		t.Fatalf("keys do not match: %v", keys)
	}

	all := registry.All()

	if len(all) != 2 {
		t.Fatalf("len mismatch: %d vs. 2", len(all))
	}

	var (
		exist bool
		repo  repository.Repository
	)

	repo, exist = all["foo"]
	if !exist || repo != repo1 {
		t.Fatal("unable to retrieve foo")
	}

	repo, exist = all["bar"]
	if !exist || repo != repo2 {
		t.Fatalf("unable to retrieve bar")
	}
}
