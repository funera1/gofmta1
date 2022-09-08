/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc/comment"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)


func GetAst(filename string) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return astFile, fset, nil
}

/*
TrimCommentMarker はコメントからコメントマーカ(// や　/*)を取り除く
pkg.go.dev/go/doc/commentによると(commemt.Parser).Parseの引数にコメントを与えるとき
コメントマーカを削除してから与えることになっているため
*/
func TrimCommentMarker(comment string) (string, string) {
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

// FormatCodeInComment はコメントを与えて、フォーマットしたコメントを返す
func FormatCodeInComment(commentString string) (string, error) {
	var p comment.Parser
	// p.Parseにつっこむときはコメントマーカー(//, /*, */)削除してから突っ込まないとだめ
	c, commentMarker := TrimCommentMarker(commentString)
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

		// cmntにはgo doc $FILEPATHの出力結果が入力されることを期待
		b, err := exec.Command("gofmt", target).Output()
		if err != nil {
			log.Fatal(err)
		}
		code := string(b)
		ast := GetAst(code)

		// コメント部分についてのみFormatCodeを適用する
		// そうしないとコメント内部でないソースコードにフォーマットがかかってしまう
		for _, cmntGrp := range ast.Comments {
			for _, cmnt := range cmntGrp.List {
				code = FormatCode(cmnt, code)

			}


		default:
			// ディレクトリ下のすべてのファイルをfilesに追加する
			var files []string
			err = filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if !d.IsDir() {
					files = append(files, path)
				}
				return err
			})
			if err != nil {
				errs = append(errs, err)
				continue
			}

			// TODO: 79行目と同じ処理なのでまとめたい
			for _, file := range files {
				// skip not gofile
				if !IsGoFile(file) {
					continue
				}
				err := GofmtalMain(file, out)
				if err != nil {
					errs = append(errs, err)
					continue
				}
			}
		}
	}
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofmtal",
	Short: "gofmtal is extended source code functionality in comments to gofmt.",
	Long:  "",
	RunE:  runE,

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		// fmt.Fprintln(os.Stderr, err)
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
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
