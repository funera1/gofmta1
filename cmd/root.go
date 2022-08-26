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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// return formatted code
func FormatCode(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
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
	if err != nil {
		return "", err
	}
	formattedCode := string(b)
	return formattedCode, nil
}

func IsGoFile(filename string) bool {
	return (filepath.Ext(filename) == ".go")
}

func GofmtalMain(filename string, writer io.Writer) error {
	formattedCode, err := FormatCode(filename)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, formattedCode)
	if err != nil {
		return err
	}

	return nil
}

func runE(cmd *cobra.Command, args []string) error {
	formatWriter := os.Stdout
	for _, arg := range args {
		switch info, err := os.Stat(arg); {
		case err != nil:
			return err
		case !info.IsDir():
			// skip not gofile
			if !IsGoFile(arg) {
				continue
			}
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
				// skip not gofile
				if !IsGoFile(file) {
					continue
				}
				err := GofmtalMain(file, formatWriter)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
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
