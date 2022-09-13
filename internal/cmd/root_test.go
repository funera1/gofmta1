package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/funera1/gofmtal/internal/derror"
	"github.com/funera1/gofmtal/internal/format"
	"github.com/google/go-cmp/cmp"
)

var (
	flagUpdate bool
)

func init() {
	derror.IsDebug = true
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func Test(t *testing.T) {
	// determine input files
	match, err := filepath.Glob("testdata/*.input")
	if err != nil {
		t.Fatal(err)
	}

	for _, in := range match {
		name := filepath.Base(in)
		t.Run(name, func(t *testing.T) {
			out := in // for files where input and output are identical
			if strings.HasSuffix(in, ".input") {
				out = in[:len(in)-len(".input")] + ".golden"
			}

			got, err := format.ProcessFile(in)
			if err != nil {
				t.Error(err)
				return
			}

			b, err := os.ReadFile(out)
			if err != nil {
				t.Error(err)
				return
			}
			want := string(b)

			if diff := cmp.Diff(got, want); diff != "" {
				t.Error(diff)
			}
		})
	}
}
