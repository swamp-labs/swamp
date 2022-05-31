package parser

import (
	"regexp"
)

const exp string = `\${((\w)+)}`

// Parser represents a parsing to apply on a field or map
type Parser interface {
	Parse(m map[string]string)
}

// Field struct defines a request field with a string
// can be initiated with the request url or other string type field
type Field struct {
	Field *string
}

// ArrayOfMap defines a field with a type []map[string]string
// can be initiated with an array of headers or query parameters
type ArrayOfMap struct {
	M *[]map[string]string
}

// Parse function search for the constant regex expression into the field
// and then replace the submatch
// key found with the m[key] value passed in argument
func (f *Field) Parse(m map[string]string) {
	re, _ := regexp.Compile(exp)
	result := re.FindAllStringSubmatch(*f.Field, -1)
	for i := range result {

		if len(result[i]) > 1 {
			key := result[i][1]
			newRegex := `\${` + key + `}`
			re, _ := regexp.Compile(newRegex)
			if m[key] != "" {
				*f.Field = re.ReplaceAllString(*f.Field, m[key])
			}
		}
	}
}

// Parse function search for the constant regex expression into
// the array of maps and then replace the submatch
// key found with the m[key] value passed in argument
func (a *ArrayOfMap) Parse(m map[string]string) {
	re, _ := regexp.Compile(exp)
	// loop over each array
	for _, n := range *a.M {
		// loop over each map
		for mapKeys, mapValue := range n {
			// Verify if there are variable replacement to do
			result := re.FindAllStringSubmatch(mapValue, -1)
			for i := range result {

				if len(result[i]) > 1 {
					key := result[i][1]
					newRegex := `\${` + key + `}`
					re, _ := regexp.Compile(newRegex)
					if m[key] != "" {
						n[mapKeys] = re.ReplaceAllString(mapValue, m[key])
					}
				}
			}
		}
	}
}
