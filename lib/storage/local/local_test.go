package local

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/DirtyHairy/csync/lib/storage"
)

func getFSInstace() storage.StorageProvider {
	fs, err := NewLocalFS("./test_artifacts")

	if err != nil {
		panic(err)
	}

	return fs
}

func getFSRoot() storage.Directory {
	root, err := getFSInstace().Root()

	if err != nil {
		panic(err)
	}

	return root
}

func getTempFSInstance() (storage.StorageProvider, error) {
	path, err := ioutil.TempDir("", "csync_test")

	if err != nil {
		return nil, err
	}

	fs, err := NewLocalFS(path)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func getTempFSRoot() (storage.Directory, error) {
	fs, err := getTempFSInstance()

	if err != nil {
		return nil, err
	}

	root, err := fs.Root()

	if err != nil {
		return nil, err
	}

	return root, nil
}

func destroyTempFS(fs storage.Directory) error {
	directory, ok := fs.(*directory)

	if !ok {
		return errors.New("no an instance of directory")
	}

	err := os.RemoveAll(directory.Entry().(*directoryEntry).realPath())

	return err
}

func checkDirectoryContents(directory storage.Directory, expectedContents []string) (entries map[string]storage.Entry, e error) {
	entries = make(map[string]storage.Entry)
	e = nil

	contents := make([]string, 0, 10)

	for entry, err := directory.NextEntry(); entry != nil; entry, err = directory.NextEntry() {
		if err != nil {
			e = errors.New(fmt.Sprintf("error while iterating over directory entries: %v", err))
			return
		}

		contents = append(contents, entry.Name())
		entries[entry.Name()] = entry
	}

	if len(contents) != len(expectedContents) {
		e = errors.New(fmt.Sprintf("expected 14 dir entries, got %d instead", len(contents)))
		return
	}

	sort.Sort(sort.StringSlice(contents))
	sort.Sort(sort.StringSlice(expectedContents))

	for idx, filename := range contents {
		if filename != expectedContents[idx] {
			e = errors.New(fmt.Sprintf("directory listing differs at %d: expected %s, got %s", idx, expectedContents[idx], filename))
			return
		}
	}

	return
}

func expectedRootEntries() []string {
	return []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "foo", "bar"}
}

func TestInstantiation(t *testing.T) {
	var (
		err error
	)

	_, err = NewLocalFS("./test_artifacts")

	if err != nil {
		t.Fatalf("unable to instantiate test FS, got %v", err)
	}

	_, err = NewLocalFS("./invalid_directory")

	if err == nil {
		t.Fatalf("should not be able to initialize FS with non-existent directory")
	}

	_, err = NewLocalFS("./local_test.go")

	if err == nil {
		t.Fatalf("should not have been able to initialize FS with file")
	}
}

func TestDirectoryListing(t *testing.T) {
	fs := getFSRoot()

	entries, err := checkDirectoryContents(fs, expectedRootEntries())

	if err != nil {
		t.Fatal(err)
	}

	var ok bool

	if _, ok = entries["a"].(storage.FileEntry); !ok {
		t.Fatalf("'a' should be a file")
	}

	if _, ok = entries["foo"].(storage.DirectoryEntry); !ok {
		t.Fatalf("foo should be a directory")
	}
}

func TestStat(t *testing.T) {
	var err error

	fs := getFSRoot()

	foo, err := fs.Stat("foo")

	if err != nil {
		t.Fatalf("failed to stat 'foo': %v", err)
	}

	a, err := fs.Stat("a")

	if err != nil {
		t.Fatalf("failed to stat 'a': %v", err)
	}

	huppe, err := fs.Stat("huppe")

	if err != nil {
		t.Fatalf("stating 'huppe' failed: %v", err)
	}

	if huppe != nil {
		t.Fatalf("stating huppe should return nil")
	}

	var ok bool

	if _, ok = a.(storage.FileEntry); !ok {
		t.Fatalf("'a' did not stat as a file")
	}

	if _, ok = foo.(storage.DirectoryEntry); !ok {
		t.Fatalf("'foo' did not stat as a directory")
	}
}

func TestNestedStat(t *testing.T) {
	fs := getFSRoot()

	_, err := fs.Stat("foo/bar")

	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
}

func TestFileNameAndPath(t *testing.T) {
	fs := getFSRoot()

	bar, _ := fs.Stat("foo/bar")

	if bar.Name() != "bar" {
		t.Fatalf("name should be 'bar', got %s", bar.Name())
	}

	if bar.Path() != "/foo/bar" {
		t.Fatalf("path should be '/foo/bar', got %s", bar.Name())
	}
}

func TestDirNameAndPath(t *testing.T) {
	fs := getFSRoot()

	hanni, _ := fs.Stat("foo/hanni")

	if hanni.Name() != "hanni" {
		t.Fatalf("name should be 'hanni'")
	}

	if hanni.Path() != "/foo/hanni" {
		t.Fatalf("path should be '/foo/hanni'")
	}
}

func TestDirRecursion(t *testing.T) {
	var err error

	fs := getFSRoot()

	foo, _ := fs.Stat("foo")

	directory, err := foo.(storage.DirectoryEntry).Open()

	if err != nil {
		t.Fatalf("failed to open directory: %v")
	}

	_, err = checkDirectoryContents(directory, []string{"bar", "hanni"})

	if err != nil {
		t.Fatal(err)
	}
}

func TestFileRead(t *testing.T) {
	var err error

	fs := getFSRoot()

	bar, _ := fs.Stat("foo/bar")

	file, err := bar.(storage.FileEntry).Open()

	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	buffer := make([]byte, 3)

	bytesRead, err := file.Read(buffer)

	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if bytesRead != 3 {
		t.Fatalf("bytes read: expected 3, got %d", bytesRead)
	}

	if string(buffer) != "baz" {
		t.Fatalf("file contents mismatch; got %v", string(buffer))
	}

	bytesRead, err = file.Read(buffer)

	if err != io.EOF {
		t.Fatalf("consecutive reads should return EOF, got %v:", err)
	}

	if bytesRead != 0 {
		t.Fatalf("bytes read: expected 0, got %d")
	}

	err = file.Close()

	if err != nil {
		t.Fatalf("close failed: %v", err)
	}

	file, err = bar.(storage.FileEntry).Open()

	if err != nil {
		t.Fatalf("second open failed: %v", err)
	}

	buffer = make([]byte, 10)

	bytesRead, err = file.Read(buffer)

	if bytesRead != 3 {
		t.Fatalf("bytes read should be 3, got %d", bytesRead)
	}

	if string(buffer[:3]) != "baz" {
		t.Fatalf("file contents differ, got: %v", string(buffer[:3]))
	}
}

func TestDirectoryRewind(t *testing.T) {
	var err error

	fs := getFSRoot()

	for i := 0; i < 5; i++ {
		_, _ = fs.NextEntry()
	}

	err = fs.Rewind()

	if err != nil {
		t.Fatalf("rewind failed: %v", err)
	}

	_, err = checkDirectoryContents(fs, expectedRootEntries())

	if err != nil {
		t.Fatal(err)
	}
}

func testCreateFile(target storage.Directory, path, referencePath string) error {
	file, err := target.CreateFile(path)

	if err != nil {
		return errors.New(fmt.Sprintf("creating file failed: %v", err))
	}

	fileClosed := false
	defer func() {
		if !fileClosed {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}
	}()

	bytesWritten, err := file.Write([]byte("Hello world"))

	if err != nil {
		return errors.New(fmt.Sprintf("failed to write to file: %v", err))
	}

	if bytesWritten != 11 {
		return errors.New(fmt.Sprintf("failed to write the full buffer, wrote %d bytes", bytesWritten))
	}

	if actualPath := file.Entry().Path(); actualPath != referencePath {
		return errors.New(fmt.Sprintf("new file has wrong path: expected %s, got %s", actualPath))
	}

	err = file.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("error closing newly created file: %v", err))
	}
	fileClosed = true

	fileEntry, err := target.Stat(path)

	if err != nil {
		return errors.New(fmt.Sprintf("failed to stat newly created file: %v", err))
	}

	if actualPath := fileEntry.Path(); actualPath != referencePath {
		return errors.New(fmt.Sprintf("new file has wrong path after statting: expected %s, got %s", actualPath))
	}

	fileRO, err := fileEntry.(storage.FileEntry).Open()

	if err != nil {
		return errors.New(fmt.Sprintf("failed to  opeen new file: %v", err))
	}

	defer func() {
		if err := fileRO.Close(); err != nil {
			panic(err)
		}
	}()

	buffer := make([]byte, 20)

	bytesRead, err := fileRO.Read(buffer)

	if err != nil && err != io.EOF {
		return errors.New(fmt.Sprintf("failed to read from newly created file: %v", err))
	}

	if bytesRead != 11 || string(buffer[:11]) != "Hello world" {
		return errors.New(fmt.Sprintf("invalid contents in written file: %v", string(buffer[:bytesRead+2])))
	}

	return nil
}

func testCreateDir(fs storage.Directory, path, referencePath string) error {
	dir, err := fs.Mkdir(path)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating directory: %v", err))
	}

	if actualPath := dir.Entry().Path(); actualPath != referencePath {
		return errors.New(fmt.Sprintf("new directory has wrong path; expected %s, got %s", referencePath, actualPath))
	}

	err = dir.Close()

	if err != nil {
		return errors.New(fmt.Sprintf("failed to close newly created directory: %v", err))
	}

	entry, err := fs.Stat(path)

	if err != nil {
		return errors.New(fmt.Sprintf("failed to stat new directory: %v", err))
	}

	if _, ok := entry.(storage.DirectoryEntry); !ok {
		return errors.New("new directory failed to stat as a directory")
	}

	if actualPath := entry.Path(); actualPath != referencePath {
		return errors.New(fmt.Sprintf("new directory has wrong path after statting; expected %s, got %s", referencePath, actualPath))
	}

	return nil
}

func TestCreateFile(t *testing.T) {
	fs, err := getTempFSRoot()

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := destroyTempFS(fs); err != nil {
			panic(err)
		}
	}()

	if err := testCreateFile(fs, "foo", "/foo"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateSingleDir(t *testing.T) {
	fs, err := getTempFSRoot()

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := destroyTempFS(fs); err != nil {
			panic(err)
		}
	}()

	if err := testCreateDir(fs, "foo", "/foo"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateNestedDir(t *testing.T) {
	fs, err := getTempFSRoot()

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := destroyTempFS(fs); err != nil {
			panic(err)
		}
	}()

	if err := testCreateDir(fs, "/foo/bar/baz", "/foo/bar/baz"); err != nil {
		t.Fatal(err)
	}

	if err := testCreateDir(fs, "/foo/bar/baz/hanni/nanni", "/foo/bar/baz/hanni/nanni"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateFileNestedPath(t *testing.T) {
	fs, err := getTempFSRoot()

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := destroyTempFS(fs); err != nil {
			panic(err)
		}
	}()

	directory, _ := fs.Mkdir("foo/bar")

	if err := testCreateFile(directory, "hanni", "/foo/bar/hanni"); err != nil {
		t.Fatal(err)
	}
}

func TestSetMtime(t *testing.T) {
	fs := getFSRoot()
	tempFs, err := getTempFSRoot()

	if err != nil {
		t.Fatal(err)
	}

	a, err := fs.Stat("/a")

	if err != nil {
		t.Fatalf("unable to stat 'a': %v", err)
	}

	aNewFile, err := tempFs.CreateFile("/a")

	if err != nil {
		t.Fatalf("unable to create 'a': %v", err)
	}

	if err := aNewFile.Close(); err != nil {
		t.Fatalf("unable to close 'a': %v", err)
	}

	aNew := aNewFile.Entry()

	if err := aNew.SetMtime(a.Mtime()); err != nil {
		t.Fatalf("failed to set mtime: %v", err)
	}

	if aNew.Mtime().Unix() != a.Mtime().Unix() {
		t.Fatalf("mtimes differ by at least one second; expected %v, got %v", a.Mtime(), aNew.Mtime())
	}
}

func TestRemoveFile(t *testing.T) {
	root, err := getTempFSRoot()

	if err != nil {
		t.Fatal(err)
	}

	file, err := root.CreateFile("hanni")

	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	if err := file.Entry().Remove(); err != nil {
		t.Fatalf("failed to remove file: %v", err)
	}

	entry, err := root.Stat("hanni")

	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	if entry != nil {
		t.Fatalf("file could not be removed")
	}
}

func TestRemoveDir(t *testing.T) {
	root, err := getTempFSRoot()

	if err != nil {
		t.Fatal(err)
	}

	directory, err := root.Mkdir("/hanni/nanni")

	if err != nil {
		t.Fatalf("unable to create dir: %v", err)
	}

	file, err := directory.CreateFile("fanni")

	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	if err := directory.Close(); err != nil {
		t.Fatalf("failed to close dir: %v", err)
	}

	entry, err := root.Stat("hanni")

	if err != nil {
		t.Fatalf("failed to stat new directory: %v", err)
	}

	if err := entry.Remove(); err != nil {
		t.Fatalf("failed to remove directory: %v", err)
	}

	entry, err = root.Stat("hanni")

	if err != nil {
		t.Fatalf("final stat failed: %v", err)
	}

	if entry != nil {
		t.Fatalf("directory could not be removed")
	}
}
