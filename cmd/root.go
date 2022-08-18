/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go/ast"
	"go/doc/comment"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/spf13/cobra"
)

func FormatCode(cmnt string) string {
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
	return string(pr.Comment(doc))
}

// とりあえず一つのファイルについてASTを返す
func GetAst(filename string) (f *ast.File, err error) {
	fset := token.NewFileSet()
	f, err = parser.ParseFile(fset, filename, nil, parser.ParseComments)
	return
}

type CodeBlock struct {
	IsComment bool
	Text      string
}

func DevideIntoCommentAndNonComment(code string, ast *ast.File) []CodeBlock {
	// コメントの位置がわかれば良さそう
	var blocks []string
	var devidePos []int = []int{0}
	isCommentMap := map[int]bool{}
	// 区切る位置のリストを取得
	for _, cmntGrp := range ast.Comments {
		for _, cmnt := range cmntGrp.List {
			pos := int(cmnt.Slash) - 1
			offset := len(cmnt.Text)
			devidePos = append(devidePos, pos)
			devidePos = append(devidePos, pos+offset)

			isCommentMap[pos] = true
		}
	}
	// 最後の位置も分割位置として含めておくことで実装が楽になる
	devidePos = append(devidePos, len([]rune(code)))

	// devidePosに従ってcodeを分割する
	for i := 0; i < len(devidePos)-1; i++ {
		start := devidePos[i]
		end := devidePos[i+1]
		blocks = append(blocks, string([]rune(code)[start:end]))
	}

	// blocks[i]がコメントかどうかの値をつける
	var codeBlocks []CodeBlock
	for i := 0; i < len(devidePos)-1; i++ {
		pos := devidePos[i]
		if _, ok := isCommentMap[pos]; ok {
			codeBlocks = append(codeBlocks, CodeBlock{
				IsComment: true,
				Text:      blocks[i],
			})
		} else {
			codeBlocks = append(codeBlocks, CodeBlock{
				IsComment: false,
				Text:      blocks[i],
			})
		}
	}
	return codeBlocks
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofmtal",
	Short: "A brief description of your application",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		ast := GetAst(target)

		// codeについてコメントとコメントでわけてブロックにする
		codeBlocks := DevideIntoCommentAndNonComment(target, ast)

		// コメント部分についてのみFormatCodeを適用する
		for i, cb := range codeBlocks {
			if cb.IsComment {
				codeBlocks[i].Text = FormatCode(cb.Text)
			}
		}

		// codeBlocksをもとに戻す
		var formattedCode string
		for _, cb := range codeBlocks {
			formattedCode += cb.Text
		}

		// TODO: 出力先を指定できるようにする
		fmt.Println(formattedCode)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gofmtal.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
