package main

import (
	"errors"
	"fmt"
)

var (
	ErrImpossibleMergeDiffTypes = errors.New("it is impossible to merge two lists together because they store a different type of data by key")
)

// mergeMaps merges two maps.
func mergeMaps(left, right map[any]any) map[any]any {
	out := make(map[any]any, len(left))
	for k, v := range left {
		out[k] = v
	}

	for k, v := range right {
		// check if v is a map
		if v, ok := v.(map[any]any); ok {
			leftval, ok := out[k]
			if !ok {
				out[k] = v
				continue
			}

			if leftval, ok := leftval.(map[any]any); ok {
				out[k] = mergeMaps(leftval, v)
				continue
			} else {
				// can't merge { "a": {} } + { "a": [1,2,3] }
				panic(fmt.Errorf("%w: %v", ErrImpossibleMergeDiffTypes, k))
			}
		}

		// check if v is slice of map
		if v, ok := v.([]any); ok {
			if _, ok := out[k]; !ok {
				out[k] = v
				continue
			}

			for i := range v {
				if leftval, ok := out[k].([]any); ok {
					out[k] = append(leftval, v[i])
				} else {
					// can't merge { "a": [1,2,3] } + { "a": {} }
					panic(fmt.Errorf("%w: %v", ErrImpossibleMergeDiffTypes, k))
				}
			}
			continue
		}

		// if another type
		out[k] = v
	}
	return out
}
