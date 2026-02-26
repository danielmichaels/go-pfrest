package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: specprep <input.json> <output.json>\n")
		os.Exit(1)
	}

	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fatal("read input: %v", err)
	}

	var spec map[string]any
	if err := json.Unmarshal(raw, &spec); err != nil {
		fatal("parse JSON: %v", err)
	}

	stats := &stats{}

	simplifyErrorSchemaData(spec, stats)
	simplifyNon200Responses(spec, stats)
	stripHTML(spec, stats)
	fixEscapedSlashes(spec, stats)

	out, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		fatal("marshal output: %v", err)
	}

	if err := os.WriteFile(os.Args[2], out, 0o644); err != nil {
		fatal("write output: %v", err)
	}

	fmt.Printf("specprep complete:\n")
	fmt.Printf("  error schema data fields simplified: %d\n", stats.dataSimplified)
	fmt.Printf("  non-200 response allOf removed:      %d\n", stats.allOfRemoved)
	fmt.Printf("  HTML tags stripped:                   %d\n", stats.htmlStripped)
	fmt.Printf("  escaped slashes fixed:               %d\n", stats.slashesFixed)
	fmt.Printf("  output: %s (%d bytes)\n", os.Args[2], len(out))
}

type stats struct {
	dataSimplified int
	allOfRemoved   int
	htmlStripped   int
	slashesFixed   int
}

var errorSchemaNames = map[string]bool{
	"AuthenticationError":      true,
	"ConflictError":            true,
	"FailedDependencyError":    true,
	"ForbiddenError":           true,
	"MediaTypeError":           true,
	"MethodNotAllowedError":    true,
	"NotAcceptableError":       true,
	"NotFoundError":            true,
	"ServerError":              true,
	"ServiceUnavailableError":  true,
	"UnprocessableContentError": true,
	"ValidationError":          true,
	"Success":                  true,
}

func simplifyErrorSchemaData(spec map[string]any, s *stats) {
	components, ok := getMap(spec, "components")
	if !ok {
		return
	}
	schemas, ok := getMap(components, "schemas")
	if !ok {
		return
	}

	for name := range errorSchemaNames {
		schema, ok := getMap(schemas, name)
		if !ok {
			continue
		}
		props, ok := getMap(schema, "properties")
		if !ok {
			continue
		}
		data, ok := getMap(props, "data")
		if !ok {
			continue
		}
		if _, hasOneOf := data["oneOf"]; hasOneOf {
			delete(data, "oneOf")
			data["type"] = "object"
			s.dataSimplified++
		}
	}
}

func simplifyNon200Responses(spec map[string]any, s *stats) {
	paths, ok := getMap(spec, "paths")
	if !ok {
		return
	}

	for _, pathItem := range paths {
		pi, ok := pathItem.(map[string]any)
		if !ok {
			continue
		}
		for _, opVal := range pi {
			op, ok := opVal.(map[string]any)
			if !ok {
				continue
			}
			responses, ok := getMap(op, "responses")
			if !ok {
				continue
			}
			for code, respVal := range responses {
				if code == "200" {
					continue
				}
				resp, ok := respVal.(map[string]any)
				if !ok {
					continue
				}
				content, ok := getMap(resp, "content")
				if !ok {
					continue
				}
				appJSON, ok := getMap(content, "application/json")
				if !ok {
					continue
				}
				schema, ok := getMap(appJSON, "schema")
				if !ok {
					continue
				}
				ref := extractErrorRef(schema)
				if ref == "" {
					continue
				}
				for k := range schema {
					delete(schema, k)
				}
				schema["$ref"] = ref
				s.allOfRemoved++
			}
		}
	}
}

func extractErrorRef(schema map[string]any) string {
	allOf, ok := schema["allOf"].([]any)
	if !ok || len(allOf) == 0 {
		return ""
	}
	for _, item := range allOf {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		ref, ok := m["$ref"].(string)
		if ok && strings.Contains(ref, "Error") {
			return ref
		}
	}
	return ""
}

func fixEscapedSlashes(spec map[string]any, s *stats) {
	walkStrings(spec, func(v string) string {
		if !strings.Contains(v, `\/`) {
			return v
		}
		n := strings.Count(v, `\/`)
		s.slashesFixed += n
		return strings.ReplaceAll(v, `\/`, `/`)
	})
}

var htmlTagRe = regexp.MustCompile(`</?(?:br|h\d?|a|strong|em|code|p|ul|li|ol|div|span|table|tr|td|th|thead|tbody)[^>]*>`)

func stripHTML(spec map[string]any, s *stats) {
	walkStrings(spec, func(v string) string {
		result := htmlTagRe.ReplaceAllStringFunc(v, func(tag string) string {
			s.htmlStripped++
			lower := strings.ToLower(tag)
			if strings.HasPrefix(lower, "<br") {
				return "\n"
			}
			return ""
		})
		return result
	})
}

func walkStrings(v any, fn func(string) string) any {
	switch val := v.(type) {
	case map[string]any:
		for k, child := range val {
			val[k] = walkStrings(child, fn)
		}
		return val
	case []any:
		for i, child := range val {
			val[i] = walkStrings(child, fn)
		}
		return val
	case string:
		return fn(val)
	default:
		return v
	}
}

func getMap(parent map[string]any, key string) (map[string]any, bool) {
	v, ok := parent[key]
	if !ok {
		return nil, false
	}
	m, ok := v.(map[string]any)
	return m, ok
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "specprep: "+format+"\n", args...)
	os.Exit(1)
}
