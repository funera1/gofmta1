package format

import (
	"bytes"
	"go/doc/comment"
	"go/format"
	"strings"
)

// 後で整理するためにprocessFileというFormatCodeの仮の関数の用意
func ProcessFile(filename string) (string, error) {
	// TODO: fはわかりにくそう
	file, err := Parse(filename)
	if err != nil {
		return "", err
	}

	// 与えられたファイルからコメントを抜き出してすべてにフォーマットをかけて戻す
	// cmnts: astからcommentGroupを抜き出したもの
	// cmnt: commentGroupからcommnetを抜き出したもの
	for i, cmnts := range file.Syntax.Comments {
		for j, cmnt := range cmnts.List {
			formattedComment, err := formatCodeInComment(cmnt.Text)
			if err != nil {
				return "", err
			}

			// フォーマットしたコメントをもとに戻す
			cmnt.Text = formattedComment
			cmnts.List[j] = cmnt
		}

		file.Syntax.Comments[i] = cmnts
	}

	var buf bytes.Buffer
	err = format.Node(&buf, file.Fset, file.Syntax)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FormatCodeInComment はコメントを与えて、フォーマットしたコメントを返す
func formatCodeInComment(commentString string) (string, error) {
	var p comment.Parser
	// p.Parseにつっこむときはコメントマーカー(//, /*, */)削除してから突っ込まないとだめ
	c, commentMarker := trimCommentMarker(commentString)
	doc := p.Parse(c)

	// commentStringからCodeを抜き出しその部分にだけフォーマットかける
	for _, c := range doc.Content {
		switch c := c.(type) {
		case *comment.Code:
			src, err := format.Source([]byte(c.Text))
			if err != nil {
				return "", err
			}
			c.Text = string(src)
		}
	}

	// コメントから抜き出したコードについてフォーマットをかける
	var pr comment.Printer
	b, err := format.Source(pr.Comment(doc))
	if err != nil {
		return "", err
	}
	formattedComment := string(b)

	// 改行するとコメントがずれるので削除
	formattedComment = strings.Trim(formattedComment, "\n")

	// コメントマーカーをつけ直す
	if commentMarker == "//" {
		formattedComment = "// " + c
	} else {
		formattedComment = "/*\n" + c + "\n*/"
	}

	return formattedComment, nil
}

/*
TrimCommentMarker はコメントからコメントマーカ(// や　/*)を取り除く
pkg.go.dev/go/doc/commentによると(commemt.Parser).Parseの引数にコメントを与えるとき
コメントマーカを削除してから与えることになっているため
*/
func trimCommentMarker(comment string) (string, string) {
	var commentMarker string
	if strings.HasPrefix(comment, "//") {
		comment = strings.TrimLeft(comment, "//")
		commentMarker = "//"
	} else {
		comment = strings.TrimLeft(comment, "/*")
		comment = strings.TrimRight(comment, "*/")
		commentMarker = "/*"
	}
	comment = strings.TrimLeft(comment, "\t")
	return comment, commentMarker
}
