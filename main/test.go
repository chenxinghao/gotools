package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var nameMap map[string]int

func getCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}

func getFileList(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		Strs := strings.Split(filepath.Base(path), ".")
		fmt.Println(Strs[0])
		if Strs[1] != "txt" {
			nameMap[Strs[0]] = 0
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func readStrLines(filePath string) ([]string, error) {

	b, err := ioutil.ReadFile(filePath)

	if err != nil {

		return nil, err

	}

	lines := strings.Split(string(b), "\n")

	return lines, nil

}

//@cxh
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	fmt.Println("running")
	return nBytes, err
}

func main() {
	nameMap = make(map[string]int)
	var fileName string
	currentDirectory := getCurrentDirectory()
	//currentDirectory:="D:\\photo"
	fmt.Scanln(&fileName)
	getFileList(currentDirectory)

	os.MkdirAll(currentDirectory+"/0", os.ModePerm)
	os.MkdirAll(currentDirectory+"/1", os.ModePerm)
	os.MkdirAll(currentDirectory+"/2", os.ModePerm)
	os.MkdirAll(currentDirectory+"/3", os.ModePerm)

	fileLines, _ := readStrLines(currentDirectory + "/" + fileName)
	for _, s := range fileLines {
		s := "C68A" + strings.Trim(strings.TrimSuffix(s, "\r"), " ")
		nameMap[s] = nameMap[s] + 1
	}

	for k, v := range nameMap {
		k := "/" + k + ".jpg"
		copy(currentDirectory+k, currentDirectory+"/"+strconv.Itoa(v)+k)
	}

}
