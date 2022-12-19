package templateString

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

const templateStringExpressionRegex string = `\${var\.((\w)+)}`
const bodyFileTemplateExpressionRegex string = `\${file:(.+)}`

type TemplateString struct {
	Format string   // Should be something like https://swamp.com/%s/%s from URL https://swamp.com/${id}/${simulation}
	Keys   []string // These are the extracted variables from URL ["id", "simulation"]
}

// UnmarshalYAML decodes a yaml node to convert it into a template string
func (t *TemplateString) UnmarshalYAML(node *yaml.Node) error {
	value := node.Value
	re, _ := regexp.Compile(bodyFileTemplateExpressionRegex)
	if re.MatchString(value) {
		filename := re.FindAllStringSubmatch(node.Value, -1)
		value = readBody(filename[0][1])
	}

	re, _ = regexp.Compile(templateStringExpressionRegex)
	matches := re.FindAllStringSubmatch(value, -1)
	t.Keys = make([]string, len(matches), cap(matches))

	for i, match := range matches {
		if len(match) > 0 {
			t.Keys[i] = match[1]
		}
	}

	if len(matches) == 0 {
		t.Keys = nil
	}
	t.Format = re.ReplaceAllString(value, "%s")
	return nil
}

// ToString returns the templateString formatted as a string
// using fmt.Sprintf function
func (t *TemplateString) ToString(context map[string]string) (string, error) {

	values := make([]string, len(t.Keys), cap(t.Keys))
	for i, key := range t.Keys {
		values[i] = context[key]
	}

	i := toInterface(values)

	return fmt.Sprintf(t.Format, i...), nil
}

func toInterface(list []string) []interface{} {
	values := make([]interface{}, len(list))
	for i, v := range list {
		values[i] = v
	}
	return values
}

func readBody(filename string) string {
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return err.Error()
	}
	body := string(byteValue)
	return body
}
