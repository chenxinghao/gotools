package AstUtils

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"gotools/Util/FileUtils"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

type Annotations struct {
	//文件地址
	FilePath            string
	CoreStr             []byte
	ImportStr           string
	FSet                *token.FileSet
	File                *ast.File
	SelectedFuncDecls   []*ast.FuncDecl
	OriginalAnnotations []string
	OppoAnnotations     []string
}

func (annotations *Annotations) InitAnnotations(filePath, importStr string) error {
	var err error
	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	annotations.FilePath = filePath
	annotations.CoreStr = code
	annotations.FSet = token.NewFileSet()
	annotations.File, err = parser.ParseFile(annotations.FSet, "", code, 4)
	if err != nil {
		return err
	}
	annotations.SelectedFuncDecls = make([]*ast.FuncDecl, 0, 20)
	annotations.ImportStr = importStr

	return nil
}

//查找被注解的方法 并删除注解
func (annotations *Annotations) SearchAnnotations(AnnotationsStr string) {
	annotations.OppoAnnotations = append(annotations.OppoAnnotations, "!"+AnnotationsStr)
	annotations.OriginalAnnotations = append(annotations.OriginalAnnotations, AnnotationsStr)
	if annotations.File.Comments != nil {
		annotations.File.Comments = annotations.File.Comments[:0]
	}
	AnnotationsStr = "//@" + AnnotationsStr
	ast.Inspect(annotations.File, func(node ast.Node) bool {
		fd, ok := node.(*ast.FuncDecl)
		if !ok {
			return true //继续搜索
		}
		if fd.Doc != nil && len(fd.Doc.List) >= 1 {
			for _, comment := range fd.Doc.List {
				text := comment.Text
				if strings.HasPrefix(text, AnnotationsStr) {
					annotations.SelectedFuncDecls = append(annotations.SelectedFuncDecls, fd)
				}
			}
			fd.Doc.List = fd.Doc.List[:0]
		}
		return false
	})

}

//TODO  一旦存在方法中不存在代码会导致生成出的代码换行有问题
func (annotations *Annotations) AnnotationsHandler(insertFuncName string, insertFuncMap map[string][]string) {
	for _, funcDecl := range annotations.SelectedFuncDecls {
		var params []string
		for _, field := range funcDecl.Type.Params.List {
			for _, name := range field.Names {
				params = append(params, name.Obj.Name)
			}
		}

		var tempStmt []ast.Stmt
		var AfterReturnFflag bool = false
		if statusList, ok := insertFuncMap[insertFuncName]; ok {
			var prefixEs *ast.ExprStmt
			var suffixEs *ast.DeferStmt
			if len(statusList) > 0 {
				as := getArgsAssignStmt(params)
				tempStmt = append(tempStmt, as)
			}
			for _, status := range statusList {

				switch status {
				case "Prefix":
					prefixEs = getInsertExprStmt(insertFuncName+"Prefix", funcDecl.Name.Name)
				case "Suffix":
					suffixEs = getInsertDeferStmt(insertFuncName+"Suffix", funcDecl.Name.Name)

				case "AfterReturn":
					AfterReturnFflag = true
				default:
				}
			}
			if prefixEs != nil {
				tempStmt = append(tempStmt, prefixEs)
			}
			if suffixEs != nil {
				tempStmt = append(tempStmt, suffixEs)
			}
			tempStmt = append(tempStmt, funcDecl.Body.List...)

			funcDecl.Body.List = tempStmt
		} else {
			return
		}

		if AfterReturnFflag {
			ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
				if node == nil {
					return false
				}
				if bs, ok := node.(*ast.BlockStmt); ok {
					for i, stmt := range bs.List {
						if _, ok1 := stmt.(*ast.ReturnStmt); ok1 {
							tempBs := []ast.Stmt{}
							tempBs = append(tempBs, bs.List[0:i]...)
							es := getInsertExprStmt(insertFuncName+"AfterReturn", funcDecl.Name.Name)
							tempBs = append(tempBs, es)
							tempBs = append(tempBs, bs.List[i:]...)
							bs.List = tempBs
						}
					}
				}
				return true
			})
		}

		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			if node == nil {
				return false
			}
			if _, ok := node.(*ast.InterfaceType); ok {
				return false
			}
			nodeValue := reflect.ValueOf(node).Elem()
			fieldNum := nodeValue.NumField()
			for i := 0; i < fieldNum; i++ {
				if _, ok := nodeValue.Field(i).Interface().(token.Pos); ok {
					nodeValue.Field(i).Set(reflect.ValueOf(token.Pos(0)))
				}
			}
			return true
		})

	}
	imports := getImportSpec(annotations.ImportStr)
	for _, decl := range annotations.File.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			d.Specs = append(d.Specs, imports)
			break
		}
	}
}

func getImportSpec(importStr string) *ast.ImportSpec {
	imports := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: "\"" + importStr + "\""}}
	return imports
}

func getArgsAssignStmt(params []string) *ast.AssignStmt {
	lhs := []ast.Expr{&ast.Ident{Name: "args", Obj: &ast.Object{}}}
	// interfacetype  必须要配置 Methods：ast.FieldList{Opening:1,Closing:2}
	// interface{}的的{}没法闭合 。导致最后ast to  string失败
	// 1和2这个是用于计算{}的偏移量，正常的情况这两个值是代码中真是偏移量
	t := &ast.ArrayType{Elt: &ast.InterfaceType{Methods: &ast.FieldList{Opening: 1, Closing: 2}}}
	var elts []ast.Expr
	//elts:=[]ast.Expr{ast.NewIdent("a"),ast.NewIdent("b"),ast.NewIdent("c")}
	for _, d := range params {
		elts = append(elts, ast.NewIdent(d))
	}
	rhs := []ast.Expr{&ast.CompositeLit{Type: t, Elts: elts}}

	as := &ast.AssignStmt{Lhs: lhs, Tok: token.DEFINE, Rhs: rhs}
	return as
}
func getInsertExprStmt(insertFuncName, originFuncName string) *ast.ExprStmt {
	//[]ast.Expr{ast.NewIdent("[]interface{}{a,b,c}")}}
	se := &ast.SelectorExpr{ast.NewIdent("AnnotationsFunc"), ast.NewIdent(insertFuncName)}
	args := []ast.Expr{&ast.Ident{Name: "args", Obj: &ast.Object{}}, ast.NewIdent("\"" + originFuncName + "\"")}
	ce := &ast.CallExpr{Fun: se, Args: args}
	es := &ast.ExprStmt{ce}
	return es
}

func getInsertDeferStmt(insertFuncName, originFuncName string) *ast.DeferStmt {

	se := &ast.SelectorExpr{ast.NewIdent("AnnotationsFunc"), ast.NewIdent(insertFuncName)}
	args := []ast.Expr{&ast.Ident{Name: "args", Obj: &ast.Object{}}, ast.NewIdent("\"" + originFuncName + "\"")}
	ce := &ast.CallExpr{Fun: se, Args: args}
	ds := &ast.DeferStmt{0, ce}
	return ds
}

func Parser(code []byte, insertFuncCountMap map[string]int) error {
	var err error
	FSet := token.NewFileSet()
	File, err := parser.ParseFile(FSet, "", code, 0)
	if err != nil {
		return err
	}

	//ast.Print(FSet,File)
	ast.Inspect(File, func(node ast.Node) bool {
		fd, ok := node.(*ast.FuncDecl)
		if !ok {
			return true //继续搜索
		}
		funcName := fd.Name.Name
		if _, ok := insertFuncCountMap[funcName]; ok {
			insertFuncCountMap[funcName] = 1
		}
		return false
	})
	return nil
}

func (annotations *Annotations) GetBuildCondition() []string {
	var fio FileUtils.FileIO
	line := fio.FindWithPrefix(annotations.CoreStr, "// +build", "package")
	if line != "" {
		line = strings.TrimPrefix(line, "// +build")
		line = strings.TrimSpace(line)
		return strings.Split(line, " ")
	}
	return nil
}

func (annotations *Annotations) GetBuildComment(flag int) string {
	var oa []string
	Conditions := annotations.GetBuildCondition()
	switch flag {
	case 1:
		oa = annotations.OppoAnnotations
	case 2:
		oa = annotations.OriginalAnnotations
	default:
		oa = annotations.OppoAnnotations
	}
	if Conditions == nil && len(oa) <= 0 {
		return ""
	}

	//TODO slice 去重
	var strMap = make(map[string]string)
	for _, v := range oa {
		strMap[v] = ""
	}
	var secondStr []string
	for key, _ := range strMap {
		secondStr = append(secondStr, key)
	}

	return strings.Join(secondStr, " ")

}

func (annotations *Annotations) SetComment() error {
	var fio FileUtils.FileIO
	var buf bytes.Buffer
	comment := annotations.GetBuildComment(1)
	if comment == "" {
		return nil
	}
	re, _ := regexp.Compile(`// \+build\s*.*\n\n`)
	coreStr := re.ReplaceAllString(string(annotations.CoreStr), "")

	buf.Write([]byte("// +build " + comment + "\n\n"))
	buf.Write([]byte(coreStr))
	_, err := fio.WriteFile(annotations.FilePath, buf.Bytes(), false)
	if err != nil {
		return err
	}
	return nil

}

//
func (annotations *Annotations) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteString("// +build " + annotations.GetBuildComment(2) + "\n\n")
	format.Node(&buf, annotations.FSet, annotations.File)
	return buf.Bytes()
}
func (annotations *Annotations) ToString() string {
	return string(annotations.ToBytes())
}
func (annotations *Annotations) PrintAST() {
	ast.Print(annotations.FSet, annotations.File)
}
