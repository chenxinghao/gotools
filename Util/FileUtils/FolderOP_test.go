package FileUtils

import (
	"testing"
)

func TestCreateALL(t *testing.T) {
	fop := &FolderOP{}
	fop.CreateALL("C:\\workspace\\GoProject\\gotools\\src\\gotools\\2231\\README.md")

}
