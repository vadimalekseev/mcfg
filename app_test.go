package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"regexp"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestFilesByPattern(t *testing.T) {
	tcs := []struct {
		Pattern string
		FS      fs.FS
		Expect  []string
	}{
		{
			Pattern: "values_production.*\\.yaml",
			FS: fstest.MapFS{
				"testdata/values_production_2.yaml": emptyDefaultFile(),
				"testdata/values_production.yaml":   emptyDefaultFile(),
				"testdata/values_stage.yaml":        emptyDefaultFile(),
				"testdata/values_stage_2.yaml":      emptyDefaultFile(),
			},
			Expect: []string{
				"testdata/values_production.yaml", "testdata/values_production_2.yaml",
			},
		},
	}

	for _, tc := range tcs {
		got := filesByPattern(tc.FS, ".", regexp.MustCompile(tc.Pattern))
		assert.Equal(t, tc.Expect, got)
	}

}

func emptyDefaultFile() *fstest.MapFile {
	return &fstest.MapFile{
		Data: []byte{},
		Mode: fs.ModePerm,
	}
}

func TestApp(t *testing.T) {
	args := []string{"merge", "-p", "^values_production.*\\.yaml", "-f", "testdata"}

	b := bytes.NewBuffer(nil)

	noError(t, run(args, b))
	got := make(map[any]any)
	noError(t, yaml.NewDecoder(b).Decode(&got))

	merged, err := os.Open("testdata/_values_production_merged.yaml")
	noError(t, err)
	expect := make(map[any]any)
	noError(t, yaml.NewDecoder(merged).Decode(&expect))

	noError(t, customCmp(expect, got))
}

// customCmp is a crutch for comparing map[any]any,
// which can have slices inside themselves, the order of which is unknown in advance.
// This func may be redudant, but now testify.Equal(Values) and reflect.DeepEqual
// can't correctly compare two maps.
func customCmp(expect, got any) error {
	switch target := expect.(type) {
	case map[any]any:
		targetb := got.(map[any]any)
		for ak, av := range target {
			bv, ok := targetb[ak]
			if !ok {
				return fmt.Errorf("not found expected key: %v", ak)
			}
			err := customCmp(av, bv)
			if err != nil {
				return err
			}
		}
	case []any:
		for _, i1 := range target {
			found := false
			for _, i2 := range got.([]any) {
				err := customCmp(i1, i2)
				if err != nil {
					continue
				}

				found = true
			}

			if !found {
				return fmt.Errorf("can't find array: %v", i1)
			}
		}
	default:
		if !reflect.DeepEqual(target, got) {
			return fmt.Errorf("%v != %v", target, got)
		}
	}

	return nil
}
