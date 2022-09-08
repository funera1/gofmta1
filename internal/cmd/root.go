/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/funera1/gofmtal/internal/format"
	"github.com/spf13/cobra"
)

const (
	ExitOK    int = 0
	ExitError int = 1
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return ExitError
	}
	return ExitOK
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofmtal",
	Short: "gofmtal is extended source code functionality in comments to gofmt.",
	Long:  "",
	RunE:  runE,
}

func runE(cmd *cobra.Command, args []string) error {
	// TODO: 自由に指定できるようにする
	var out io.Writer
	out = os.Stdout

	var errs []error

	// argがファイルかディレクトリかそれ以外かで場合分け
	for _, arg := range args {
		switch info, err := os.Stat(arg); {

		case err != nil:
			errs = append(errs, err)
			log.Println("not file or dir")
			continue

		case !info.IsDir():
			err := GofmtalMain(arg, out)
			if err != nil {
				log.Println("miss GofmtalMain")
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
				log.Println("miss filepath.WalkDir")
				errs = append(errs, err)
				continue
			}

			for _, file := range files {
				isGofile := bool(filepath.Ext(file) == ".go")
				if isGofile {
					continue
				}

				err := GofmtalMain(file, out)
				if err != nil {
					log.Println("miss GofmtalMain")
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

func GofmtalMain(filename string, writer io.Writer) error {
	// formattedCode, err := processFile(filename)
	formattedCode, err := format.ProcessFile(filename)
	if err != nil {
		log.Println("miss format.ProcessFile")
		return err
	}

	_, err = fmt.Fprintln(writer, formattedCode)
	if err != nil {
		log.Println("miss fmt.Fprintln")
		return err
	}

	return nil
}
