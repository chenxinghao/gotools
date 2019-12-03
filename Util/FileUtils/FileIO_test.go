package FileUtils

import (
	"fmt"
	"testing"
)

func TestReadStrLines(t *testing.T) {
	fio := &FileIO{}
	ss, _ := fio.ReadStrLines("C:\\workspace\\GoProject\\gotools\\src\\gotools\\README.md")
	for _, s := range ss {
		fmt.Println(s)
	}
}

func TestReadByteLines(t *testing.T) {
	fio := &FileIO{}
	ss, _ := fio.ReadByteLines("C:\\workspace\\GoProject\\gotools\\src\\gotools\\README.md")
	for _, s := range ss {
		fmt.Println(string(s))
	}
}
func TestUpdateFileByLine(t *testing.T) {
	fio := &FileIO{}
	fio.UpdateFileByLine("C:\\workspace\\GoProject\\gotools\\src\\gotools\\README.md", []byte("完善中ing"), 2)
}
