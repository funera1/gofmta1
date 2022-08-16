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

func FormatCode(cmnt *ast.Comment, code string) string {
	// コメントのソースコードに対しフォーマットをかける
	var p comment.Parser
	doc := p.Parse(cmnt.Text)
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
	format_cmnt := string(pr.Comment(doc))
	rune_code := []rune(code)
	// もとのコメントをフォーマットかけたもので置き換える
	// 実行するとわかるが、もとのコメントのフォーマットしたコメントの長さが変わるのでこれでは
	// 置き換えれていない
	for i, nc := range format_cmnt {
		rune_code[int(cmnt.Slash)-1+i] = nc
	}
	return string(rune_code)
}

func GetAst(code string) *ast.File {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		log.Fatalln("Error", err)
	}
	return f
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

		// コメント部分についてのみFormatCodeを適用する
		// そうしないとコメント内部でないソースコードにフォーマットがかかってしまう
		for _, cmntGrp := range ast.Comments {
			for _, cmnt := range cmntGrp.List {
				code = FormatCode(cmnt, code)
			}
		}

		print(code)
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
