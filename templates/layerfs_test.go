package templates

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
)

func TestLayerFs_Stat(t *testing.T) {
	lfs := LayerFs{}
	if err := testStat(lfs, true, 0); err != nil {
		t.Error(err)
	}
	lfs = append(lfs, os.DirFS("../templates/included"))
	if err := testStat(lfs, true, 0); err != nil {
		t.Error(err)
	}

}

func testStat(lfs LayerFs, expectDir bool, expectSize int) error {
	fi, err := fs.Stat(lfs, ".")
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("Expected root of layerfs to be a directory")
	}
	if fi.Size() != 0 {
		return fmt.Errorf("Expected zero size for root of layer fs, got %d", fi.Size())
	}
	return nil
}
