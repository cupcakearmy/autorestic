package internal

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestOptionToStringWithoutPrefix(t *testing.T) {
	opt := "test"
	result := optionToString(opt)
	assertEqual(t, result, "--test")
}

func TestOptionsToStringWithSinglePrefix(t *testing.T) {
	opt := "-test"
	result := optionToString(opt)
	assertEqual(t, result, "-test")
}

func TestOptionsToStringWithDoublePrefix(t *testing.T) {
	opt := "--test"
	result := optionToString(opt)
	assertEqual(t, result, "--test")
}

func TestAppendOneOptionToSlice(t *testing.T) {
	result := []string{}
	optionMap := OptionMap{"string-flag": []interface{}{"/root"}}

	appendOptionsToSlice(&result, optionMap)
	expected := []string{
		"--string-flag", "/root",
	}
	assertSliceEqual(t, result, expected)
}

func TestAppendBoolOptionToSlice(t *testing.T) {
	result := []string{}
	optionMap := OptionMap{"boolean-flag": []interface{}{true}}

	appendOptionsToSlice(&result, optionMap)
	expected := []string{
		"--boolean-flag",
	}
	assertSliceEqual(t, result, expected)
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

func assertEqual(t testing.TB, result, expected string) {
	t.Helper()

	if result != expected {
		t.Errorf("got %s, want %s", result, expected)
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
