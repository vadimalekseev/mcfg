package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"regexp"
	"sort"

	"gopkg.in/yaml.v2"
)

var (
	ErrInvalidCommand = errors.New("invalid command: example usage: merge -p values_production.*\\.yaml")
)

func run(rawArgs []string, out io.Writer) error {
	if len(rawArgs) < 1 {
		return ErrInvalidCommand
	}

	switch rawArgs[0] {
	case "merge":
		// pass
	default:
		return ErrInvalidCommand
	}

	args, err := parseArgs(rawArgs[1:])
	if err != nil {
		return err
	}

	compiledPattern, err := regexp.Compile(args.Pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	openPath := path.Join(currentDir, args.Folder)
	targetFolderFiles := filesByPattern(os.DirFS(openPath), ".", compiledPattern)

	sort.Strings(targetFolderFiles)

	if len(targetFolderFiles) == 0 {
		return fmt.Errorf("can't find any file")
	}

	result, err := walkAndMerge(openPath, targetFolderFiles)
	if err != nil {
		return err
	}

	yaml.NewEncoder(out).Encode(result)

	return nil
}

type Args struct {
	Pattern string
	Folder  string
}

func parseArgs(args []string) (Args, error) {
	f := flag.NewFlagSet("", flag.ExitOnError)

	pattern := f.String("p", "YAML files pattern", "files pattern")
	folder := f.String("f", ".", "search folder")

	if err := f.Parse(args); err != nil {
		return Args{}, err
	}

	return Args{
		Pattern: *pattern,
		Folder:  *folder,
	}, nil
}

// filesByPattern returns all matched files in fs.FS
func filesByPattern(target fs.FS, root string, pattern *regexp.Regexp) []string {
	var result []string
	_ = fs.WalkDir(target, root, func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		name := d.Name()
		if pattern.Match([]byte(name)) {
			result = append(result, path)
		}
		return nil
	})

	return result
}

// walkAndMerge merges files into one map.
func walkAndMerge(basePath string, targetFolderFiles []string) (map[any]any, error) {
	result := make(map[any]any)
	for _, subpath := range targetFolderFiles {
		f, err := os.Open(path.Join(basePath, subpath))
		if err != nil {
			return nil, fmt.Errorf("can't open file %w", err)
		}

		current := make(map[any]any)
		yaml.NewDecoder(f).Decode(&current)

		if err := f.Close(); err != nil {
			return nil, fmt.Errorf("can't close file: %w", err)
		}

		result = mergeMaps(result, current)
	}

	return result, nil
}
