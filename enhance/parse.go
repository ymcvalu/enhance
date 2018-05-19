package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func ParseFile(filepath string, file interface{}) *File {
	_file := &File{}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath, file, parser.ParseComments)
	_file.ImportPos = getLine(fset, f.Package) + 1
	_file.ImportKey = true
	//ast.Print(fset, f)
	if err != nil {
		panic(err)
	}
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				if d.Lparen == token.NoPos {
					_file.ImportPos = getLine(fset, d.TokPos) + 1
					_file.ImportKey = true
				} else {
					_file.ImportPos = getLine(fset, d.Rparen)
					_file.ImportKey = false
				}

			}
		case *ast.FuncDecl:
			//ast.Print(fset, fn)
			//if has the doc
			if d.Doc == nil {
				continue
			}

			funName := d.Name.Name
			if strings.HasPrefix(funName, "_") {
				continue
			}
			annos, exists := ParseAnnos(fset, d.Doc)
			if !exists {
				continue
			}

			fun := Func{}
			fun.Annos = annos

			fun.Name = funName

			fun.Rbrace = fset.Position(d.Body.Rbrace).Line
			fun.Lbrace = fset.Position(d.Body.Lbrace).Line

			if ins := d.Type.Params; ins != nil {
				fun.Ins = ParseParams(fset, ins.List)
			}

			if outs := d.Type.Results; outs != nil {
				fun.Outs = ParseParams(fset, outs.List)
			}

			if recv := d.Recv; recv != nil {
				fun.Recv = parseRecv(fset, recv.List[0])
			}
			_file.Funs = append(_file.Funs, fun)
		}

	}
	return _file
}

func getLine(fset *token.FileSet, pos token.Pos) int {
	return fset.Position(pos).Line
}
func ParseParams(fset *token.FileSet, list []*ast.Field) []Param {
	params := make([]Param, 0)
	for _, in := range list {
		pTyp := parseTyp(in.Type)
		param := Param{
			Type: pTyp,
		}
		for _, n := range in.Names {
			param.Names = append(param.Names, n.Name)
		}
		params = append(params, param)
	}
	return params
}

func parseRecv(fset *token.FileSet, rf *ast.Field) Recv {
	recv := Recv{}
	recv.Type = parseTyp(rf.Type)
	if len(rf.Names) > 0 {
		recv.Name = rf.Names[0].Name
	} else {
		recv.Name = "_self_"
	}
	return recv
}

func parseTyp(typ ast.Expr) string {
	typName := ""
	switch t := typ.(type) {
	case *ast.Ident:
		typName = t.Name
	case *ast.ChanType:
		switch t.Dir {
		case 1:
			typName = "<-chan "
		case 2:
			typName = "chan<- "
		default:
			typName = "chan "
		}
		typName += parseTyp(t.Value)
	case *ast.MapType:
		typName = "map[" + parseTyp(t.Key) + "]" + parseTyp(t.Value)
	case *ast.ArrayType:
		if t.Len != nil {
			typName = fmt.Sprintf("[%s]", t.Len.(*ast.BasicLit).Value)
		} else {
			typName = "[]"
		}
		typName += parseTyp(t.Elt)
	case *ast.SelectorExpr:
		typName = t.X.(*ast.Ident).Name + "." + t.Sel.Name
	case *ast.StarExpr:
		typName = "*" + parseTyp(t.X)
	}
	return typName
}

const TAG = "enhance:"

func ParseAnnos(fset *token.FileSet, cg *ast.CommentGroup) ([]Anno, bool) {
	annos := make([]Anno, 0)
	for _, c := range cg.List {
		comment := c.Text
		if !strings.HasPrefix(comment, "//") {
			continue
		}
		//remove the `//`
		comment = strings.Trim(comment[2:], " ")
		if !strings.HasPrefix(comment, TAG) {
			continue
		}
		comment = strings.Trim(comment[len(TAG):], " ")
		if comment == "" {
			continue
		}
		tags := strings.Split(comment, ",")
		for i := range tags {
			tags[i] = fmt.Sprintf(`"%s"`, strings.Trim(tags[i], " "))
		}
		anno := Anno{
			Tags: tags,
			Line: getLine(fset, c.Pos()),
		}
		annos = append(annos, anno)
	}

	return annos, len(annos) > 0
}
