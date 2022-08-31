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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// for debug
func DbgFset(fset *token.FileSet) {
	println("DEBUG: fset")
	fset.Iterate(
		func(f *token.File) bool {
			if f == nil {
				return false
			}
			fmt.Printf("name is %s\n", f.Name())
			fmt.Printf("base is %d\n", f.Base())
			fmt.Printf("size is %d\n", f.Size())
			return true
		},
	)
}

// for debug
func DbgComments(cmnts *ast.CommentGroup) {
	println("DEBUG: cmnts")
	for _, cmnt := range cmnts.List {
		fmt.Println(cmnt.Text)
	}
}

func GetAst(filename string) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return f, fset, nil
}

// 後で整理するためにprocessFileというFormatCodeの仮の関数の用意
func processFile(filename string) error {
	// TODO: fはわかりにくそう
	f, fset, err := GetAst(filename)
	if err != nil {
		return err
	}

	// TODO: めちゃくちゃややこしいのでわかりやすくする
	// cmnts: astからcommentGroupを抜き出したもの
	// cmnt: commentGroupからcommnetを抜き出したもの
	var p comment.Parser
	for i, cmnts := range f.Comments {
		for j, cmnt := range cmnts.List {
			fmt.Printf("%d, %d: \n", i, j)
			fmt.Println("before format: " + cmnt.Text)

			// p.Parseにつっこむときはコメントマーカー(//, /*, */)削除してから突っ込まないとだめ
			// TODO: 関数化する
			c := strings.TrimLeft(cmnt.Text, "/*")
			c = strings.TrimRight(c, "*/")
			c = strings.Trim(c, "\t")
			println(c)
			doc := p.Parse(c)

			// cmntからCodeを抜き出しその部分にだけフォーマットかける
			for _, c := range doc.Content {
				switch c := c.(type) {
				case *comment.Code:
					src, err := format.Source([]byte(c.Text))
					if err != nil {
						return err
					}
					c.Text = string(src)
				}
			}

			var pr comment.Printer
			b, err := format.Source(pr.Comment(doc))
			if err != nil {
				return err
			}

			// TODO: tabの個数とかの調整をする必要がある
			c = string(b)
			c = strings.Trim(c, "\t")
			c = strings.Trim(c, "\n")
			// TODO: もとのコメントマーカーを覚えておいてそれに戻す
			c = "/*\n" + c + "\n*/"
			fmt.Println(c)
			// :EYE:
			cmnt.Text = c
			cmnts.List[j] = cmnt
		}

		f.Comments[i] = cmnts
	}

	ast.Print(fset, f)
	// TODO: 多分fsetが原因だが、出力するときにコメントがちょっとずれる
	format.Node(os.Stdout, fset, f)
	return nil
}

// return formatted code
// func FormatCode(filename string) (string, error) {
// 	b, err := os.ReadFile(filename)
// 	if err != nil {
// 		return "", err
// 	}
// 	code := string(b)
//
// 	// コード中のソースコードを抜き出しフォーマットをかける
// 	var p comment.Parser
// 	doc := p.Parse(code)
// 	for _, c := range doc.Content {
// 		switch c := c.(type) {
// 		case *comment.Code:
// 			src, err := format.Source([]byte(c.Text))
// 			if err != nil {
// 				return "", err
// 			}
// 			c.Text = string(src)
// 		}
// 	}
//
// 	// コード全体に対してフォーマットをかける
// 	var pr comment.Printer
// 	b, err = format.Source(pr.Comment(doc))
// 	if err != nil {
// 		return "", err
// 	}
// 	formattedCode := string(b)
// 	return formattedCode, nil
// }

func IsGoFile(filename string) bool {
	return (filepath.Ext(filename) == ".go")
}

func GofmtalMain(filename string, writer io.Writer) error {
	// formattedCode, err := processFile(filename)
	err := processFile(filename)
	if err != nil {
		return err
	}

	// _, err = fmt.Fprintln(writer, formattedCode)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func runE(cmd *cobra.Command, args []string) error {
	// TODO: 自由に指定できるようにする
	var out io.Writer
	out = os.Stdout

	var errs []error

	for _, arg := range args {
		switch info, err := os.Stat(arg); {

		case err != nil:
			errs = append(errs, err)
			continue

		case !info.IsDir():
			// skip not gofile
			if !IsGoFile(arg) {
				continue
			}
			err := GofmtalMain(arg, out)
			if err != nil {
				errs = append(errs, err)
				continue
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
