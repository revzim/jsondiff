package jsondiff

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var testNameReplacer = strings.NewReplacer(",", "", "(", "", ")", "")

type testcase struct {
	Name   string      `json:"name"`
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
	Patch  []Operation `json:"patch"`
}

func TestArrayCases(t *testing.T)  { runCasesFromFile(t, "testdata/tests/array.json") }
func TestObjectCases(t *testing.T) { runCasesFromFile(t, "testdata/tests/object.json") }
func TestRootCases(t *testing.T)   { runCasesFromFile(t, "testdata/tests/root.json") }

func TestOptions(t *testing.T) {
	makeopts := func(opts ...Option) []Option { return opts }

	for _, tt := range []struct {
		testfile string
		options  []Option
	}{
		{"testdata/tests/options/invertible.json", makeopts(Invertible())},
		{"testdata/tests/options/factorization.json", makeopts(Factorize())},
		{"testdata/tests/options/rationalization.json", makeopts(Rationalize())},
		{"testdata/tests/options/all.json", makeopts(Factorize(), Rationalize(), Invertible())},
	} {
		var (
			ext  = filepath.Ext(tt.testfile)
			base = filepath.Base(tt.testfile)
			name = strings.TrimSuffix(base, ext)
		)
		t.Run(name, func(t *testing.T) {
			runCasesFromFile(t, tt.testfile, tt.options...)
		})
	}
}

func runCasesFromFile(t *testing.T, filename string, opts ...Option) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var cases []testcase
	if err := json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	runTestCases(t, cases, opts...)
}

func runTestCases(t *testing.T, cases []testcase, opts ...Option) {
	for _, tc := range cases {
		name := testNameReplacer.Replace(tc.Name)

		t.Run(name, func(t *testing.T) {
			d := differ{}
			d.applyOpts(opts...)
			d.diff(tc.Before, tc.After)

			if d.patch != nil {
				t.Logf("\n%s", d.patch)
			}
			if len(d.patch) != len(tc.Patch) {
				t.Errorf("got %d patches, want %d", len(d.patch), len(tc.Patch))
				return
			}
			for i, op := range d.patch {
				want := tc.Patch[i]
				if g, w := op.Type, want.Type; g != w {
					t.Errorf("op #%d mismatch: op: got %q, want %q", i, g, w)
				}
				if g, w := op.Field.String(), want.Field.String(); g != w {
					t.Errorf("op #%d mismatch: field: got %q, want %q", i, g, w)
				}
				switch want.Type {
				case OperationCopy, OperationMove:
					if g, w := op.From.String(), want.From.String(); g != w {
						t.Errorf("op #%d mismatch: from: got %q, want %q", i, g, w)
					}
				case OperationAdd, OperationReplace:
					if !reflect.DeepEqual(op.Value, want.Value) {
						t.Errorf("op #%d mismatch: value: unequal", i)
					}
				}
			}
		})
	}
}
