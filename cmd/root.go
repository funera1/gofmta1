/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go/doc/comment"
	"go/format"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// コードに対しフォーマットを掛けた文字列を返す
func FormatCode(filename string) (formattedCode string, err error) {
	// ファイルの中身を取得
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	code := string(b)

	// コード中のソースコードを抜き出しフォーマットをかける
	var p comment.Parser
	doc := p.Parse(code)
	for _, c := range doc.Content {
		switch c := c.(type) {
		case *comment.Code:
			src, err := format.Source([]byte(c.Text))
			if err == nil {
				c.Text = string(src)
			}
		}
	}

	// コード全体に対してフォーマットをかける
	var pr comment.Printer
	b, err = format.Source(pr.Comment(doc))
	formattedCode = string(b)
	return
}

func IsGoFile(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

func GofmtalMain(targetfile string, writer io.Writer) (err error) {
	// Gofileでないならフォーマットかけない
	if !IsGoFile(targetfile) {
		return
	}
	formattedCode, err := FormatCode(targetfile)
	if err != nil {
		return
	}
	fmt.Fprintln(writer, formattedCode)
	return
}

func runE(cmd *cobra.Command, args []string) (err error) {
	formatWriter := os.Stdout
	for _, arg := range args {
		switch info, err := os.Stat(arg); {
		case err != nil:
			return err
		case !info.IsDir():
			GofmtalMain(arg, formatWriter)
		default:
			// ディレクトリ下のすべてのファイルをfilesに追加する
			var files []string
			err = filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if !d.IsDir() {
					files = append(files, path)
				}
				return err
			})
			for _, file := range files {
				GofmtalMain(file, formatWriter)
			}
		}
	}
	return
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofmtal",
	Short: "gofmtal is extended source code functionality in comments to gofmt.",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runE,
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
