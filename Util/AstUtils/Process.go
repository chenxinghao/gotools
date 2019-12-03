package AstUtils

import (
	"fmt"
	"gotools/Util/FileUtils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//通过文件路径和文函数名来获取第一个匹配的内容和文件信息
func CheckFuncWithRes(path string, insertFuncCountMap map[string]int) {
	err := filepath.Walk(path, func(filename string, fi os.FileInfo, err error) error {
		name := fi.Name()
		if !strings.HasSuffix(name, ".go") {
			return nil
		}
		Core, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if e := Parser(Core, insertFuncCountMap); e == nil {
			return filepath.SkipDir
		}

		return nil

	})
	if err != nil {
		panic(err)
	}

}

func CheckFunc(path string, insertFuncMap map[string][]string) {
	insertFuncCountMap := make(map[string]int)
	var prefix int
	var suffix int
	var afterReturn int
	for k, _ := range insertFuncMap {
		insertFuncCountMap[k+"Prefix"] = 0
		insertFuncCountMap[k+"Suffix"] = 0
		insertFuncCountMap[k+"AfterReturn"] = 0
	}
	CheckFuncWithRes(path, insertFuncCountMap)
	for k, _ := range insertFuncMap {
		prefix = insertFuncCountMap[k+"Prefix"]
		suffix = insertFuncCountMap[k+"Suffix"]
		afterReturn = insertFuncCountMap[k+"AfterReturn"]
		if prefix == 1 {
			insertFuncMap[k] = append(insertFuncMap[k], "Prefix")
		}
		if suffix == 1 {
			insertFuncMap[k] = append(insertFuncMap[k], "Suffix")
		}
		if afterReturn == 1 {
			insertFuncMap[k] = append(insertFuncMap[k], "AfterReturn")
		}
	}
}

func Handler(filePath, annoName, importStr, funcName string, insertFuncMap map[string][]string) *Annotations {
	a := Annotations{}
	a.InitAnnotations(filePath, importStr)
	a.SearchAnnotations(annoName)
	if len(a.SelectedFuncDecls) <= 0 {
		return nil
	}
	a.PrintAST()
	a.AnnotationsHandler(funcName, insertFuncMap)
	a.SetComment()
	fmt.Print(a.ToString())
	return &a
}

func WalkAndHandler(dirPath, annoName, importStr, funcName string, insertFuncMap map[string][]string) {
	var fio FileUtils.FileIO
	err := filepath.Walk(dirPath, func(filePath string, fi os.FileInfo, err error) error {
		if fi.IsDir() { // 忽略目录
			if fi.Name() == "vendor" || fi.Name() == "AnnotationsFunc" {
				return filepath.SkipDir
			}
			return nil
		}

		name := fi.Name()
		if !strings.HasSuffix(name, ".go") {
			return nil
		}
		namePix := strings.TrimSuffix(name, ".go")

		a := Handler(filePath, annoName, importStr, funcName, insertFuncMap)
		if a != nil {
			fio.WriteFile(fmt.Sprintf("%s\\%s_Pro.go", filepath.Dir(filePath), namePix), a.ToBytes(), false)
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
