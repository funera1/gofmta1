// 使い方: go doc $FILEPATH | go run main.go
//
// 機能: go docで出力したコメント内部のソースコードに対してフォーマットをかける
package main

import (
	"go/doc/comment"
	"go/format"
	"io/ioutil"
	"os"
)

func main() {

    // cmntにはgo doc $FILEPATHの出力結果が入力されることを期待する
	var cmnt string
	bytes, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        panic(err)
    }
    cmnt = string(bytes)

    // コメントのソースコードに対しフォーマットをかける
	var p comment.Parser
	doc := p.Parse(cmnt)
	for _, c := range doc.Content {
		switch c := c.(type) {
		case *comment.Code:
			src, err := format.Source([]byte(c.Text))
			if err == nil {
				c.Text = string(src)
			}
		}
	}
	var pr comment.Printer
	os.Stdout.Write(pr.Comment(doc))
}
