package templates

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

const testDirName = "testdir"
const testFileName = "testfile.txt"
const testFileContent = "hello world\n"

func TestNewFileBlobStore(t *testing.T) {
	_, err := createStore()
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	defer func() {
		os.RemoveAll(testDirName)
	}()
}

func TestFileBlobStore_Write(t *testing.T) {
	fs, err := createStore()
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	defer func() {
		os.RemoveAll(testDirName)
	}()
	if err = fs.Write("testfile.txt", []byte(testFileContent)); err != nil {
		t.Errorf("Failed to write test file %v", err)
	}
}

func TestFileBlobStore_Read(t *testing.T) {
	fs, err := createStore()
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	defer func() {
		os.RemoveAll(testDirName)
	}()
	if err = fs.Write(testFileName, []byte(testFileContent)); err != nil {
		t.Errorf("Failed to write test file %v", err)
	}

	by, err := fs.Read(testFileName)
	if err != nil {
		t.Errorf("failed to read test file  %v", err)
	}
	if !bytes.Equal(by, []byte(testFileContent)) {
		t.Errorf("unexpected file content read from blob file.")
	}
}

func TestFileBlobStore_Contains(t *testing.T) {
	fs, err := createStore()
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	defer func() {
		os.RemoveAll(testDirName)
	}()
	if err = fs.Write(testFileName, []byte(testFileContent)); err != nil {
		t.Errorf("Failed to write test file %v", err)
	}

	if fs.Contains("unknown") {
		t.Errorf("unexpected true response from Contains, given unknown name")
	}
	if !fs.Contains(testFileName) {
		t.Errorf("unexpected false response from Contains, given known name")
	}
}

func createStore() (BlobStore, error) {
	err := os.Mkdir(testDirName, 0750)
	if err != nil {
		return nil, fmt.Errorf("failed to creare test directory  %v", err)
	}

	return NewFileBlobStore(testDirName)
}
