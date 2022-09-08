package format

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type File struct {
	Syntax *ast.File
	Fset   *token.FileSet
}

func Parse(filename string) (*File, error) {
	fset := token.NewFileSet()
	syntax, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return *File{syntax, fset}, nil
}
