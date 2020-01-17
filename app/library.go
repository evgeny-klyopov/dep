package app

import (
	"reflect"
	"regexp"
)

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func RegexSplit(text string, separator string) []string {
	regex := regexp.MustCompile(separator)
	indexes := regex.FindAllStringIndex(text, -1)
	result := make([]string, len(indexes) + 1)
	last := 0

	for i, element := range indexes {
		result[i] = text[last:element[0]]
		last = element[1]
	}
	result[len(indexes)] = text[last:len(text)]

	return result
}