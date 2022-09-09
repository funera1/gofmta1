package format

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/funera1/gofmtal/internal/derror"
)

type File struct {
	Syntax *ast.File
	Fset   *token.FileSet
}

func Parse(filename string) (_ *File, rerr error) {
	defer derror.Wrap(&rerr, "Parse(%q)", filename)

	fset := token.NewFileSet()
	syntax, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	file := &File{
		Syntax: syntax,
		Fset:   fset,
	}
	return file, nil
}
