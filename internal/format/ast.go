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
	Tfile  *token.File
	Lines  []int
}

func Parse(filename string) (_ *File, rerr error) {
	defer derror.Wrap(&rerr, "Parse(%q)", filename)

	fset := token.NewFileSet()
	syntax, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// formattedCodeはコメントがずれてる場合があるので、その部分の整形を行う
	// できるだけ処理の途途でやるべきだが、ひとまずわかりやすい場所に書く
	tfile := fset.File(syntax.Pos())
	lines := make([]int, tfile.LineCount())
	for i := 0; i < tfile.LineCount(); i++ {
		// i+1にしてるのは引数でとる行番号が1から始まるから
		lines[i] = tfile.Offset(tfile.LineStart(i + 1))
	}

	file := &File{
		Syntax: syntax,
		Fset:   fset,
		Tfile:  tfile,
		Lines:  lines,
	}
	return file, nil
}
