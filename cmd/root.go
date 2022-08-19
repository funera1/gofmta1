/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go/doc/comment"
	"go/format"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// コードに対しフォーマットを掛けた文字列を返す
func FormatCode(filename string) (formattedCode string, err error) {
	// ファイルの中身を取得
	b, err := ioutil.ReadFile(filename)
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
	var pr comment.Printer
	b, err = format.Source(pr.Comment(doc))
	formattedCode = string(b)
	return
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofmtal",
	Short: "A brief description of your application",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// TODO: 複数ファイルやディレクトリ指定したときの対応
		targetfile := args[0]
		if err != nil {
			return
		}

		formattedCode, err := FormatCode(targetfile)
		fmt.Println(formattedCode)
		return
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
