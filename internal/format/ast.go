package format

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

type File struct {
	Syntax *ast.File
	Fset   *token.FileSet
}

func Parse(filename string) (*File, error) {
	fset := token.NewFileSet()
	syntax, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Println("miss parse.ParseFile")
		return nil, err
	}

	file := &File{
		Syntax: syntax,
		Fset:   fset,
	}
	return file, nil
}
