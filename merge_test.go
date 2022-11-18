package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestYamlMerge(t *testing.T) {
	t.Run("production", func(t *testing.T) {
		first, second, third := parseYaml("testdata/values_production_monitoring.yaml"),
			parseYaml("testdata/values_production_shipping_fee.yaml"),
			parseYaml("testdata/values_production.yaml")

		result := mergeMaps(mergeMaps(first, second), third)

		expected := parseYaml("testdata/_values_production_merged.yaml")

		resultMarshalled, err := yaml.Marshal(result)
		noError(t, err)
		expectMarshalled, err := yaml.Marshal(expected)
		noError(t, err)

		a, b := string(resultMarshalled), string(expectMarshalled)
		_, _ = a, b

		equal(t, expected, result)
	})

	// should panic due to different datatypes by key
	t.Run("panic", func(t *testing.T) {
		require.PanicsWithError(t, ErrImpossibleMergeDiffTypes.Error()+": a", func() {
			first, second := parseYaml("testdata/panic_case1_1.yaml"),
				parseYaml("testdata/panic_case1_2.yaml")

			_ = mergeMaps(first, second)
		})

		require.PanicsWithError(t, ErrImpossibleMergeDiffTypes.Error()+": b", func() {
			first, second := parseYaml("testdata/panic_case2_1.yaml"),
				parseYaml("testdata/panic_case2_2.yaml")

			_ = mergeMaps(first, second)
		})
	})

	tcs := []struct {
		Name         string
		YAML1, YAML2 string
		Expected     string
	}{
		{
			Name:     "half prod",
			YAML1:    "testdata/half_prod_1.yaml",
			YAML2:    "testdata/half_prod_2.yaml",
			Expected: "testdata/_half_prod_merged.yaml",
		},
		{
			Name:     "simple",
			YAML1:    "testdata/simple_1.yaml",
			YAML2:    "testdata/simple_2.yaml",
			Expected: "testdata/_simple_merged.yaml",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			first, second := parseYaml(tc.YAML1),
				parseYaml(tc.YAML2)

			got := mergeMaps(first, second)

			expected := parseYaml(tc.Expected)

			equal(t, expected, got)
		})
	}
}

func parseYaml(path string) (result map[any]any) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	err = yaml.NewDecoder(f).Decode(&result)
	if err != nil {
		panic(err)
	}

	return result
}

func noError(t *testing.T, err error) {
	require.NoError(t, err)
}

func equal(t *testing.T, expected, got any) {
	require.Equal(t, expected, got)
}
