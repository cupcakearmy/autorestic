package options

import (
	"fmt"
	"strings"
)

type OptionMap map[string][]interface{}
type Options map[string]OptionMap

func (o Options) GetOptions(keys []string) []string {
	var selected []string
	for _, key := range keys {
		o[key].AppendOptionsToSlice(&selected)
	}
	return selected
}

func (m OptionMap) AppendOptionsToSlice(str *[]string) {
	for key, values := range m {
		for _, value := range values {
			// Bool
			asBool, ok := value.(bool)
			if ok && asBool {
				*str = append(*str, optionToString(key))
				continue
			}
			*str = append(*str, optionToString(key), fmt.Sprint(value))
		}
	}
}

func optionToString(option string) string {
	if !strings.HasPrefix(option, "-") {
		return "--" + option
	}
	return option
}
