/*
Copyright © 2022 funera1
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/funera1/gofmtal/internal/derror"
	"github.com/funera1/gofmtal/internal/file"
	"github.com/funera1/gofmtal/internal/format"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"go.uber.org/multierr"
)

const (
	ExitOK    int = 0
	ExitError int = 1
)

var (
	writeFlag *bool
)

func init() {
	// var writeFlag *bool = flag.Bool("w", false, "write result to (source) file instead of stdout")
	writeFlag = rootCmd.Flags().BoolP("write", "w", false, "write result to (source) file instead of stdout")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, rootCmd.Use+":", err)
		return ExitError
	}
	return ExitOK
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofmtal",
	Short: "gofmtal is extended source code functionality in comments to gofmt.",
	Long:  "",
	RunE:  main,
}

func main(cmd *cobra.Command, args []string) error {
	var rerr error

	// 引数でとるflag
	flags := cmd.Flags()

	// argがファイルかディレクトリかそれ以外かで場合分け
	for _, arg := range args {
		switch info, err := os.Stat(arg); {

		// not file or dir
		case err != nil:
			rerr = multierr.Append(rerr, err)
			continue

		// file
		case !info.IsDir():
			err := gofmtalMain(flags, arg, info)
			if err != nil {
				rerr = multierr.Append(rerr, err)
				continue
			}

		// dir
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
				rerr = multierr.Append(rerr, err)
				continue
			}

			for _, file := range files {
				isGofile := bool(filepath.Ext(file) == ".go")
				if isGofile {
					continue
				}

				err := gofmtalMain(flags, file, info)
				if err != nil {
					rerr = multierr.Append(rerr, err)
					continue
				}
			}
		}
	}

	if rerr != nil {
		return rerr
	}
	return nil
}

func gofmtalMain(flags *flag.FlagSet, filename string, info fs.FileInfo) (rerr error) {
	defer derror.Wrap(&rerr, "GofmtalMain(%q)", filename)

	formattedCode, err := format.ProcessFile(filename)
	if err != nil {
		return err
	}

	// もとのファイルに上書きする
	if *writeFlag {
		if info == nil {
			return fmt.Errorf("-w should not have been allowed with stdin")
		}

		// もとのファイルのbackupを取る
		perm := info.Mode().Perm()
		src, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		bakname, err := file.BackupFile(filename+".", src, perm)
		if err != nil {
			return err
		}

		// 上書きする
		err = os.WriteFile(filename, []byte(formattedCode), perm)
		if err != nil {
			// 失敗したらもとに戻す
			os.Rename(bakname, filename)
			return err
		}

		// backup fileを削除する
		err = os.Remove(bakname)
		if err != nil {
			return err
		}
	} else {
		// 標準出力
		_, err = fmt.Fprintln(os.Stdout, formattedCode)
		if err != nil {
			return err
		}
	}

	return nil
}
