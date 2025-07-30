package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/pretty"

	"github.com/wsva/lib_go/stack"
)

func Unescape(content string) string {
	content = strings.Replace(content, "\\u003c", "<", -1)
	content = strings.Replace(content, "\\u003e", ">", -1)
	content = strings.Replace(content, "\\u0026", "&", -1)
	return content
}

func Indent(jsonBytes []byte) []byte {
	return pretty.Pretty(jsonBytes)
}

func IndentString(jsonString string) string {
	return string(pretty.Pretty([]byte(jsonString)))
}

/*
format each object in json-array to single line
then easy for line processing

BUG: what if element in list hat \n?
*/
func FormatSingleLine(jsonString string) string {
	list, err := SplitString(jsonString)
	if err != nil {
		return jsonString
	}
	resultString := "[\n"
	for i := 0; i < len(list)-1; i++ {
		resultString += "  " + list[i] + ",\n"
	}
	resultString += "  " + list[len(list)-1] + "\n"
	resultString += "]"
	return resultString
}

/*
split json bytes
[{"name": "a"},	{"name": "b"}]
split to list type
*/
func Split(src []byte) ([][]byte, error) {
	var list []interface{}
	err := json.Unmarshal(src, &list)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to list error: %v", err)
	}

	var result [][]byte
	for _, v := range list {
		b, _ := json.Marshal(v)
		result = append(result, b)
	}

	return result, nil
}

func SplitString(src string) ([]string, error) {
	list, err := Split([]byte(src))
	if err != nil {
		return nil, err
	}
	var result []string
	for _, v := range list {
		result = append(result, string(v))
	}
	return result, nil
}

/*
use SplitString instead
*/
func Split2(jsonString string) ([]string, error) {
	reg := regexp.MustCompile(`^\[({.*},)*{.*}]$`)
	if !reg.MatchString(jsonString) {
		return nil, errors.New("not JSON array of object")
	}
	regBegin := regexp.MustCompile(`^\[`)
	regEnd := regexp.MustCompile(`]$`)
	jsonString = regBegin.ReplaceAllString(jsonString, "")
	jsonString = regEnd.ReplaceAllString(jsonString, "")
	var objectString string
	var objectStringList []string
	braceStack := &stack.Stack{}
	for _, v := range jsonString {
		if v == ',' && braceStack.Empty() {
			objectStringList = append(objectStringList, objectString)
			objectString = ""
			continue
		}
		objectString += string(v)
		if v == '{' {
			braceStack.Push("{")
		}
		if v == '}' {
			err := braceStack.Pop()
			if err != nil {
				return nil, err
			}
		}
	}
	objectStringList = append(objectStringList, objectString)
	return objectStringList, nil
}
