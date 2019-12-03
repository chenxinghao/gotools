package FileUtils

import (
	"os"
	"path/filepath"
	"strings"
)

type FolderOP struct {
}

func (this *FolderOP) CreateALL(filePath string) error {
	err := os.MkdirAll(filePath, os.ModePerm)
	return err

}
func (this *FolderOP) GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}
func (this *FolderOP) GetParentDirectory(dirctory string) string {
	length := strings.LastIndex(dirctory, "/")
	runes := []rune(dirctory)
	return string(runes[:length])
}
