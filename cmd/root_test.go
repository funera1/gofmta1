package cmd

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func runTest(t *testing.T, in, out string) {
	_, err := os.Lstat(in)
	if err != nil {
		t.Error(err)
		return
	}
	// got := processFile(in)
	want, err := os.ReadFile(out)
	if err != nil {
		t.Error(err)
		return
	}

	var buf bytes.Buffer
	if got := buf.Bytes(); !bytes.Equal(got, want) {
		// if *update {
		// 	if in != out {
		// 		if err := os.WriteFile(out, got, 0666); err != nil {
		// 			t.Error(err)
		// 		}
		// 		return
		// 	}
		// in == out: don't accidentally destroy input
		// t.Errorf("WARNING: -update did not rewrite input file %s", in)
		// }

		// t.Errorf("(gofmt %s) != %s (see %s.gofmt)\n%s", in, out, in,
		// 	diff.Diff("expected", want, "got", got))
		if err := os.WriteFile(in+".gofmt", got, 0666); err != nil {
			t.Error(err)
		}
	}
}

func Test(t *testing.T) {
	// determine input files
	match, err := filepath.Glob("../testdata/*.input")
	if err != nil {
		t.Fatal(err)
	}

	// add larger examples
	match = append(match, "root.go", "root_test.go")

	for _, in := range match {
		name := filepath.Base(in)
		t.Run(name, func(t *testing.T) {
			out := in // for files where input and output are identical
			if strings.HasSuffix(in, ".input") {
				out = in[:len(in)-len(".input")] + ".golden"
			}
			runTest(t, in, out)
			if in != out && !t.Failed() {
				// Check idempotence.
				runTest(t, out, out)
			}
		})
	}
}
