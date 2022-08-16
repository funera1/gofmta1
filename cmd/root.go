/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go/ast"
	"go/doc/comment"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"

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

func GetAst(code string) *ast.File {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		log.Fatalln("Error", err)
	}
	return f
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
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MinimumNArgs(1),
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		// cmntにはgo fmt targetの出力結果が入力されることを期待
		b, err := exec.Command("gofmt", target).Output()
		if err != nil {
			log.Fatal(err)
		}
		code := string(b)
		ast := GetAst(code)

		// codeについてコメントとコメントでわけてブロックにする
		codeBlocks := DevideIntoCommentAndNonComment(code, ast)

		// コメント部分についてのみFormatCodeを適用する
		for i, cb := range codeBlocks {
			if cb.IsComment {
				codeBlocks[i].Text = FormatCode(cb.Text)
			}
		}

		// codeBlocksをもとに戻す
		format_code := ""
		for _, cb := range codeBlocks {
			format_code += cb.Text
		}

		print(format_code)
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
