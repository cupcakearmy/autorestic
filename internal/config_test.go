package internal

import (
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestOptionToString(t *testing.T) {
	t.Run("no prefix", func(t *testing.T) {
		opt := "test"
		result := optionToString(opt)
		assertEqual(t, result, "--test")
	})

	t.Run("single prefix", func(t *testing.T) {
		opt := "-test"
		result := optionToString(opt)
		assertEqual(t, result, "-test")
	})

	t.Run("double prefix", func(t *testing.T) {
		opt := "--test"
		result := optionToString(opt)
		assertEqual(t, result, "--test")
	})
}

func TestAppendOneOptionToSlice(t *testing.T) {
	t.Run("string flag", func(t *testing.T) {
		result := []string{}
		optionMap := OptionMap{"string-flag": []interface{}{"/root"}}

		appendOptionsToSlice(&result, optionMap)
		expected := []string{
			"--string-flag", "/root",
		}
		assertSliceEqual(t, result, expected)
	})

	t.Run("bool flag", func(t *testing.T) {
		result := []string{}
		optionMap := OptionMap{"boolean-flag": []interface{}{true}}

		appendOptionsToSlice(&result, optionMap)
		expected := []string{
			"--boolean-flag",
		}
		assertSliceEqual(t, result, expected)
	})

	t.Run("int flag", func(t *testing.T) {
		result := []string{}
		optionMap := OptionMap{"int-flag": []interface{}{123}}

		appendOptionsToSlice(&result, optionMap)
		expected := []string{
			"--int-flag", "123",
		}
		assertSliceEqual(t, result, expected)
	})
}

func TestAppendMultipleOptionsToSlice(t *testing.T) {
	result := []string{}
	optionMap := OptionMap{
		"string-flag": []interface{}{"/root"},
		"int-flag":    []interface{}{123},
	}

	appendOptionsToSlice(&result, optionMap)
	expected := []string{
		"--string-flag", "/root",
		"--int-flag", "123",
	}
	if len(result) != len(expected) {
		t.Errorf("got length %d, want length %d", len(result), len(expected))
	}

	// checks that expected option comes after flag, regardless of key order in map
	for i, v := range expected {
		v = strings.TrimPrefix(v, "--")

		if value, ok := optionMap[v]; ok {
			if val, ok := value[0].(int); ok {
				if expected[i+1] != strconv.Itoa(val) {
					t.Errorf("Flags and options order are mismatched. got %v, want %v", result, expected)
				}
			}
		}
	}
}

func TestAppendOptionWithMultipleValuesToSlice(t *testing.T) {
	result := []string{}
	optionMap := OptionMap{
		"string-flag": []interface{}{"/root", "/bin"},
	}

	appendOptionsToSlice(&result, optionMap)
	expected := []string{
		"--string-flag", "/root",
		"--string-flag", "/bin",
	}
	assertSliceEqual(t, result, expected)
}

func TestGetOptionsOneKey(t *testing.T) {
	optionMap := OptionMap{
		"string-flag": []interface{}{"/root"},
	}
	options := Options{"backend": optionMap}
	keys := []string{"backend"}

	result := getOptions(options, keys)
	expected := []string{
		"--string-flag", "/root",
	}
	assertSliceEqual(t, result, expected)
}

func TestGetOptionsMultipleKeys(t *testing.T) {
	firstOptionMap := OptionMap{
		"string-flag": []interface{}{"/root"},
	}
	secondOptionMap := OptionMap{
		"boolean-flag": []interface{}{true},
		"int-flag":     []interface{}{123},
	}
	options := Options{
		"all":    firstOptionMap,
		"forget": secondOptionMap,
	}
	keys := []string{"all", "forget"}

	result := getOptions(options, keys)
	expected := []string{
		"--string-flag", "/root",
		"--boolean-flag",
		"--int-flag", "123",
	}
	reflect.DeepEqual(result, expected)
}

func TestSaveConfigProducesReadableConfig(t *testing.T) {
	workDir := t.TempDir()
	viper.SetConfigFile(path.Join(workDir, ".autorestic.yml"))

	// Required to appease the config reader
	viper.Set("version", 2)

	c := Config{
		Version: "2",
		Locations: map[string]Location{
			"test": {
				Type: "local",
				name: "test",
				From: []string{"in-dir"},
				To:   []string{"test"},
				// ForgetOption & ConfigOption have previously marshalled in a way that
				// can't get read correctly
				ForgetOption: "foo",
				CopyOption:   map[string][]string{"foo": {"bar"}},
			},
		},
		Backends: map[string]Backend{
			"test": {
				name: "test",
				Type: "local",
				Path: "backup-target",
				Key:  "supersecret",
			},
		},
	}

	err := c.SaveConfig()
	assert.NoError(t, err)

	// Ensure we the config reading logic actually runs
	config = nil
	once = sync.Once{}
	readConfig := GetConfig()
	assert.NotNil(t, readConfig)
	assert.Equal(t, c, *readConfig)
}

func assertEqual[T comparable](t testing.TB, result, expected T) {
	t.Helper()

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func assertSliceEqual(t testing.TB, result, expected []string) {
	t.Helper()

	if len(result) != len(expected) {
		t.Errorf("got length %d, want length %d", len(result), len(expected))
	}

	for i := range result {
		assertEqual(t, result[i], expected[i])
	}
}
