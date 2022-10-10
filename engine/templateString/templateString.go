package templateString

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

const templateStringExpressionRegex string = `\${var\.((\w)+)}`

type TemplateString struct {
	Format string   // Should be something like https://swamp.com/%s/%s from URL https://swamp.com/${id}/${simulation}
	Keys   []string // These are the extracted variables from URL ["id", "simulation"]
}

// UnmarshalYAML decodes a yaml node to convert it into a template string
func (t *TemplateString) UnmarshalYAML(node *yaml.Node) error {
	re, _ := regexp.Compile(templateStringExpressionRegex)

	matches := re.FindAllStringSubmatch(node.Value, -1)
	t.Keys = make([]string, len(matches), cap(matches))

	for i, match := range matches {
		if len(match) > 0 {
			t.Keys[i] = match[1]
		}
	}

	if len(matches) == 0 {
		t.Keys = nil
	}
	t.Format = re.ReplaceAllString(node.Value, "%s")
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
