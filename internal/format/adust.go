package format

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/exp/slices"
)

// やりたいこと
// 行の情報を増やせばフォーマットのずれは直るのか　-> 直る

// func adjustLines(file *File) {
// 	// 上記実験を踏まえて、format後のlinesを計算によって求める
// 	// var file *ast.File
// 	// var tfile *token.File
// 	// var lines

// 	// newlines := getNewLines(file.Syntax, file.Tfile, file.Lines)

// 	// debug
// 	// oldlinesとnewlinesを比較
// 	fmt.Println("oldlines:")
// 	fmt.Println(file.Lines)
// 	fmt.Println("newlines:")
// 	// fmt.Println(newlines)
// 	fmt.Println()

// 	// newlinesを適用
// 	// file.Tfile.SetLines(newlines)
// }

// ずれたlinesの調整をすべてのコメントについて行う
// func getNewLines(file *ast.File, tfile *token.File, oldlines []int) []int {
// 	lines := oldlines
// 	for _, cmnts := range file.Comments {
// 		for _, cmnt := range cmnts.List {
// 			// 行頭のリストに対応するコメントの範囲を取得
// 			startInd, endInd, err := getCommentSection(tfile, lines, cmnt)
// 			if err != nil {
// 				// startIndとendIndは行頭にあることを仮定してあるが
// 				// 当然すべてのコメントが行頭にあるわけではないので、
// 				// その条件を満たさない場合はcontinuesする
// 				continue
// 			}
//
// 			// コメントの範囲を新しいコメントで置換
// 			commentlines := getCommentLines(tfile, nil, cmnt)
// 			lines = replaceSlice(lines, startInd, endInd, commentlines)
// 		}
// 	}
// 	return lines
// }

// コメント部分の書き換えについて、linesの調整つきで書き換えを行う
func updateComment(cmnt *ast.Comment, formattedComment string, file *File) ([]int, error) {
	// 行頭のリストに対応するコメントの範囲を取得
	startInd, endInd, err := getCommentSection(file.Tfile, file.Lines, cmnt)
	if err != nil {
		return nil, err
	}

	// コメントの範囲を新しいコメントで置換
	startOfs := file.Tfile.Offset(cmnt.Pos())
	commentlines := getCommentLines(file.Tfile, startOfs, formattedComment)
	newlines := replaceSlice(file.Lines, startInd, endInd, commentlines)

	// debug
	fmt.Printf("startInd: %d, endInd: %d\n", startInd, endInd)
	fmt.Println("commentlines:")
	fmt.Println(commentlines)

	return newlines, nil
}

// コメントのoldlinesにおける対応する範囲を返す
// 対応するコメントの範囲は[startInd, endInd)である
func getCommentSection(tfile *token.File, oldlines []int, cmnt *ast.Comment) (int, int, error) {
	// コメントのoffsetを取る
	startOfs := tfile.Offset(cmnt.Pos())
	endOfs := tfile.Offset(cmnt.End())

	// offsetから対応するlinesのindexを取る
	startInd := slices.Index(oldlines, startOfs)

	// endIndは改行の数で数えたほうがいいかもしれない
	endInd := startInd + 1

	// startIndから改行がある毎に1増やしていく
	for _, c := range cmnt.Text {
		if c == '\n' {
			endInd++
		}
	}

	// startIndとendIndに値がなければエラーとする
	if startInd == -1 || endInd == -1 {
		// debug
		fmt.Println("DEBUG")
		fmt.Printf("startOfs: %d, endOfs: %d\n", startOfs, endOfs)
		fmt.Println("oldlines:")
		fmt.Println(oldlines)
		return startInd, endInd, fmt.Errorf("No value for startInd or endInd\nstartInd: %d, endInd: %d\n", startInd, endInd)
	}

	return startInd, endInd, nil
}

// コメントの行頭のオフセットを取得
func getCommentLines(tfile *token.File, startOfs int, formattedComment string) []int {
	// a. 更新されたコメントの行頭のオフセットを取って、commentlinesとする
	var commentlines []int

	// コメントの行頭のオフセットを取得
	commentlines = append(commentlines, startOfs)
	for i := 0; i < len(formattedComment); i++ {
		if formattedComment[i] == '\n' {
			// 行頭は改行の次の位置
			commentlines = append(commentlines, startOfs+i+1)
		}
	}

	return commentlines
}

// oldの[start, end)をtargetで置換する
func replaceSlice[S ~[]E, E any](old S, start, end int, target S) S {
	new := make(S, len(old)-(end-start)+len(target))
	n := copy(new, old[:start])
	n += copy(new[n:], target)
	copy(new[n:], old[end:])

	return new
}
