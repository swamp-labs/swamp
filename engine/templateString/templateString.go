package templateString

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

const templateStringExpressionRegex string = `\${((\w)+)}`

type TemplateString struct {
	Format string   // Should be something like https://swamp.com/%s/%s from URL https://swamp.com/${id}/${simulation}
	Keys   []string // These are the extracted variables from URL ["id", "simulation"]
}

func YamlNodeToTemplateString(node *yaml.Node) TemplateString {

	ts := TemplateString{}
	re, _ := regexp.Compile(templateStringExpressionRegex)

	matches := re.FindAllStringSubmatch(node.Value, -1)
	ts.Keys = make([]string, len(matches), cap(matches))

	for i, match := range matches {
		if len(match) > 0 {
			ts.Keys[i] = match[1]
		}
	}

	if len(matches) == 0 {
		ts.Keys = nil
	}
	ts.Format = re.ReplaceAllString(node.Value, "%s")

	return ts
}

func (t *TemplateString) ToString(context map[string]string) (string, error) {

	if len(t.Keys) > 0 {
		values := make([]string, len(t.Keys), cap(t.Keys))
		for i, key := range t.Keys {
			values[i] = context[key]
		}

		i := toInterface(values)

		return fmt.Sprintf(t.Format, i...), nil
	}
	return t.Format, nil
}

func toInterface(list []string) []interface{} {
	values := make([]interface{}, len(list))
	for i, v := range list {
		values[i] = v
	}
	return values
}
