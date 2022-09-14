package format

import (
	"bytes"
	"go/ast"
	"go/doc/comment"
	"go/format"
	"go/token"
	"strings"

	"github.com/funera1/gofmtal/internal/derror"
)

/* TODO:
gofmtでprocessFileという名前がつけられてたから同じ名前つけていたが、
関数名から意味を読み取りにくいので、renameしても良さそう
*/
// 後で整理するためにprocessFileというFormatCodeの仮の関数の用意
func ProcessFile(filename string) (_ string, rerr error) {
	defer derror.Wrap(&rerr, "ProcessFile(%q)", filename)

	file, err := Parse(filename)
	if err != nil {
		return "", err
	}

	// 与えられたファイルからコメントを抜き出してすべてにフォーマットをかけて戻す
	// cmnts: astからcommentGroupを抜き出したもの
	// cmnt: commentGroupからcommnetを抜き出したもの
	for i, cmnts := range file.Syntax.Comments {
		for j, cmnt := range cmnts.List {
			formattedComment, err := formatCodeInComment(cmnt, file)
			if err != nil {
				return "", err
			}

			// フォーマットしたコメントをもとに戻す
			cmnt.Text = formattedComment
			cmnts.List[j] = cmnt
		}

		file.Syntax.Comments[i] = cmnts
	}

	// formatでずれたlinesを調整する
	file.Tfile.SetLines(file.Lines)

	var buf bytes.Buffer
	err = format.Node(&buf, file.Fset, file.Syntax)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// FormatCodeInComment はコメントを与えて、フォーマットしたコメントを返す
func formatCodeInComment(cmnt *ast.Comment, file *File) (_ string, rerr error) {
	defer derror.Debug(&rerr, "formatCodeInComment(%q)", cmnt.Text)

	commentString := cmnt.Text

	var p comment.Parser
	// p.Parseにつっこむときはコメントマーカー(//, /*, */)削除してから突っ込まないとだめ
	commentInfo := trimCommentMarker(commentString)
	doc := p.Parse(commentInfo.Comment)

	// commentStringからCodeを抜き出しその部分にだけフォーマットかける
	for _, c := range doc.Content {
		switch c := c.(type) {
		case *comment.Code:
			src, err := format.Source([]byte(c.Text))
			if err != nil {
				// format.Source()でsyntax errorが発生するコードは
				// そもそもformatできないので、無視する
				continue
			}

			c.Text = string(src)
		}
	}

	var pr comment.Printer
	b := pr.Comment(doc)
	formattedComment := string(b)

	// TODO
	// formattedCommentともとのコメントの改行場所をあわせたい。そこがコメントマーカがずれる原因になってる
	// 改行するとコメントがずれるので削除
	formattedComment = strings.Trim(formattedComment, "\n")

	// コメントマーカーをつけ直す
	if commentInfo.CommentMarker == "//" {
		formattedComment = "// " + formattedComment
	} else {
		if commentInfo.LineCount == 1 {
			formattedComment = "/*" + formattedComment + "*/"
		} else {
			formattedComment = "/*\n" + formattedComment + "\n*/"
		}
	}

	// formatするとコメント内の改行が変動するので、調整する必要がある
	adjustLines(formattedComment, cmnt, file)

	return formattedComment, nil
}

/*
TrimCommentMarker はコメントからコメントマーカ(// や　/*)を取り除く
pkg.go.dev/go/doc/comment によると(commemt.Parser).Parseの引数にコメントを与えるとき
コメントマーカを削除してから与えることになっているため
*/

type CommentInfo struct {
	Comment       string
	CommentMarker string
	LineCount     int
}

func trimCommentMarker(comment string) CommentInfo {
	// 行数数える
	lineCount := strings.Count(comment, "\n")

	var commentMarker string

	// commentからcommentMarkerを取り除く
	if strings.HasPrefix(comment, "//") {
		commentMarker = "//"

		comment = strings.TrimLeft(comment, "//")
	} else {
		commentMarker = "/*"

		comment = strings.TrimLeft(comment, "/*")
		comment = strings.TrimRight(comment, "*/")
	}

	comment = strings.TrimLeft(comment, "\t")
	return CommentInfo{
		Comment:       comment,
		CommentMarker: commentMarker,
		LineCount:     lineCount,
	}
}

// file.Linesについて、cmntの範囲内のLinesについてずれた分の調整をする
func adjustLines(formattedComment string, cmnt *ast.Comment, file *File) {
	startPos := cmnt.Slash
	startOfs := file.Tfile.Offset(startPos)
	// TODO: len(cmnt.Text)ってエスケープ文字含含んでのかな？正しいpos返してる？
	endPos := token.Pos(int(startPos) + len(cmnt.Text))
	endOfs := file.Tfile.Offset(endPos)

	// TODO: startInd, endIndは境界条件について注意する
	// startIndex, endIndexはコメントマーカーを含まない
	startIndex := -1
	for i := 0; i < len(file.Lines); i++ {
		if startIndex != -1 {
			break
		}

		// startIndexはコメントが始まってから一番最初の改行の位置を持つ
		if startOfs < file.Lines[i] {
			startIndex = i
		}
	}

	endIndex := -1
	for i := 0; i < len(file.Lines)-1; i++ {
		if endIndex != -1 {
			break
		}

		// endIndexはコメントが終わってから一番最初の改行の位置を持つ
		// これはコメントの位置を配列について[startIndex, endIndex)を範囲として取るため
		if endOfs <= file.Lines[i+1] {
			endIndex = i
		}
	}

	var newlines []int
	// formattedCommentの改行のPosを取得する
	for i, c := range formattedComment {
		// 改改の位置を調べる. '\n'は10
		if c == 10 {
			// TODO: startOfs+iであってる？
			newlines = append(newlines, startOfs+i)
		}
	}

	// newlinesを[startIndex, endindex]の部分に置換する
	lis := make([][]int, 3)
	lis[0] = file.Lines[:startIndex]
	lis[1] = newlines
	lis[2] = file.Lines[endIndex+1:]

	file.Lines = myappend(lis)
}

// 配列内の配列について連結したものを返す
func myappend(lis [][]int) []int {
	var ret []int

	if len(lis) == 1 {
		ret = append(ret, lis[0]...)
	} else {
		ret = append(lis[0], myappend(lis[1:])...)
	}

	return ret
}
