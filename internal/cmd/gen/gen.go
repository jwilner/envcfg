package main

import (
	"bytes"
	"go/format"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

type option struct {
	Name, Type, Default, ZeroVal string
}

func (o option) Comma() bool {
	return o.Name == "Comma"
}

func (o option) Trait() string {
	if strings.HasSuffix(o.Name, "e") {
		return o.Name + "r"
	}
	return o.Name + "er"
}

type specCfg struct {
	MethodName, TypeName string
	ParseFunc            string
	imports              []string
	Options              []option
}

func (s specCfg) ParserName() string {
	return s.MethodName + "Parser"
}

func (s specCfg) OptName() string {
	return s.MethodName + "Opt"
}

func (s specCfg) AllOptions() []option {
	options := s.Fields()
	options = append(options, option{"Default", s.TypeName, "", ""})
	sort.Slice(options, func(i, j int) bool {
		return options[i].Name < options[j].Name
	})
	return options
}

func (s specCfg) Fields() []option {
	var os []option
	for _, o := range s.Options {
		if o.Name != "" {
			os = append(os, o)
		}
	}
	if s.Slice() {
		os = append(os, comma)
	}
	return os
}

func (s specCfg) Slice() bool {
	return strings.Contains(strings.ToLower(s.MethodName), "slice")
}

func (s specCfg) Imports() []string {
	imports := append([]string(nil), s.imports...)
	if s.Slice() && s.ParseFunc != "" {
		imports = append(imports, "fmt")
	}
	if s.Slice() {
		imports = append(imports, "encoding/json")
	}
	sort.Strings(imports)
	return imports
}

var (
	base    = option{"Base", "int", "", "0"}
	bitSize = option{"BitSize", "int", "", "0"}
	layout  = option{"Layout", "string", "time.RFC3339", `""`}
	comma   = option{"Comma", "rune", "", "0"}

	parseInt = "strconv.ParseInt"

	types = []specCfg{
		{"Bool", "bool", "strconv.ParseBool", []string{"strconv"}, []option{{}}},
		{"Duration", "time.Duration", "time.ParseDuration", []string{"time"}, []option{{}}},
		{"Float", "float64", "strconv.ParseFloat", []string{"strconv"}, []option{{}, bitSize}},
		{"Int", "int64", parseInt, []string{"strconv"}, []option{{}, base, bitSize}},
		{"IntSlice", "[]int64", parseInt, []string{"strconv"}, []option{{}, base, bitSize}},
		{"String", "string", "", nil, nil},
		{"StringSlice", "[]string", "", nil, nil},
		{"Time", "time.Time", "time.Parse", []string{"time"}, []option{layout, {}}},
		{"Uint", "uint64", "strconv.ParseUint", []string{"strconv"}, []option{{}, base, bitSize}},
	}

	uniOptTmpl = tmplWithFuncs(`
// Code generated by internal/cmd/gen/gen.go DO NOT EDIT.
package envcfg

type UniOpt interface {
	modify(s *spec) error
{{ range $t := . -}}
	modify{{ .ParserName }}(p *{{ .ParserName | unexported }}) error 
{{ end }}
}

type uniOptFunc func(s *spec) error

func (f uniOptFunc) modify(s *spec) error {
	return f(s)
}

{{ range $t := . }}
func (uniOptFunc) modify{{ .ParserName }}(p *{{ .ParserName | unexported }}) error {
	return nil
} 
{{ end }}

var _ UniOpt = new(uniOptFunc)
`)

	tmpl = tmplWithFuncs(`
// Code generated by internal/cmd/gen/gen.go DO NOT EDIT.
package envcfg

{{ if .Imports -}}
import (
{{ range $i := .Imports -}}
     "{{ . }}"
{{ end }}
)
{{ end }}

type {{ .OptName }} interface {
	modify{{ .ParserName }}(p *{{ .ParserName | unexported }}) error
}

// {{ .MethodName }} extracts and parses the variable provided according to the options provided.
// Available options:
{{ range $o := .AllOptions -}}
// - {{ .Name | snake_case }}
{{ end -}}
func (c *Cfg) {{ .MethodName }}(docOpts string, opts ...{{ .OptName }}) {{ .TypeName }} {
    s, err := new{{ .MethodName }}Spec(docOpts, opts)
	if err != nil {
		if c.panic {
			panic(err)
		}
		c.addError(err)
    }
	c.addDescription(s.describe())
	v, _ := c.evaluate(s).({{ .TypeName }})
	return v
}

var {{ .MethodName }} = struct {
	Default func({{ .TypeName }}) {{ .OptName }}
{{ range $o := .Fields -}}
	{{ .Name }} func({{ .Type }}) {{ $.OptName }}
{{ end }}
}{
	Default: func(def {{ .TypeName }}) {{ .OptName }} {
		return defaultOpt(def)
	},
{{ range $o := .Fields -}}
	{{ .Name }}: func({{ .Name | unexported }} {{ .Type }}) {{ $.OptName }} {
		return {{ $.OptName | unexported }}Func(func(p *{{ $.ParserName | unexported }}) error {
			p.set{{ .Name }}({{ .Name | unexported }})
			return nil
		})
	},
{{ end }}
}

type {{ .OptName | unexported }}Func func(p *{{ .ParserName | unexported }}) error

func (f {{ .OptName | unexported }}Func) modify{{ .ParserName }}(p *{{ .ParserName | unexported }}) error {
	return f(p)
}

var _ {{ .OptName }} = new({{ .OptName | unexported }}Func)

func new{{ .MethodName }}Spec(docOpts string, opts []{{ .OptName }}) (*spec, error) {
	parsed, err := parseDocOpts(docOpts)
	if err != nil {
		return nil, err
	}

	os := make([]{{ .MethodName }}Opt, 0, len(opts) + len(parsed))
	for _, p := range parsed {
		os = append(os, p)
	}

	p := new({{ .ParserName | unexported }})
    s := &spec{
		parser: p,
		typeName: "{{ .TypeName }}",
	}

    for _, opt := range append(os, opts...) {
		if opt, ok := opt.(interface { modify(*spec) error }); ok {
			if err = opt.modify(s); err != nil {
				return nil, err
			}
			continue
		}
		if err := opt.modify{{ .ParserName }}(p); err != nil {
			return nil, err
		}
    }
	
	if s.flags & flagDefaultValString > 0 {
		if s.defaultVal, err = p.parse(s.defaultValS); err != nil {
			return nil, err
		}
	}

	return s, nil
}

type {{ .ParserName | unexported }} struct {
{{ range $o := .Fields -}}
{{ if .Comma -}}
	slicer
{{ else -}}
	{{ .Trait | unexported }}
{{ end -}}
{{ end -}}
}

func (p *{{ .ParserName | unexported }}) parse(s string) (interface{}, error) {
{{ if .Slice -}}
{{ if not .ParseFunc -}}
	return p.parseSlice(s)
{{ else -}}
	ses, err := p.parseSlice(s)
	if err != nil {
		return nil, err
	}
	{{ range $o := .Options -}}
	{{ if .Default }}
	{{ .Name | unexported }} := p.{{ .Name | unexported }}
	if {{ .Name | unexported }} == {{ .ZeroVal }} {
		{{ .Name | unexported }} = {{ .Default }}
	}
	{{ end }}
	{{ end -}}
	vals := make({{ .TypeName }}, len(ses))
	for i, v := range ses {
		el, err := {{ .ParseFunc }}(
			{{ range $o := .Options -}}
			{{ if .Name }}{{ if not .Default }}p.{{ end }}{{ .Name | unexported }}{{ else }}v{{ end }},
			{{ end }}
		)
		if err != nil {
			return nil, fmt.Errorf("%v index: %v", i, err)
		}
		vals[i] = el
	}
	return vals, nil
{{ end }}
{{ else -}}
{{ if not .ParseFunc -}}
	return s, nil
{{ else -}}
	{{ range $o := .Options -}}
	{{ if .Default -}}
	{{ .Name | unexported }} := p.{{ .Name | unexported }}
	if {{ .Name | unexported }} == {{ .ZeroVal }} {
		{{ .Name | unexported }} = {{ .Default }}
	}
	{{ end -}}
	{{ end -}}
	return {{ .ParseFunc }}(
{{ range $o := .Options -}}
		{{ if .Name -}}{{ if not .Default }}p.{{ end }}{{ .Name | unexported }}{{ else -}}s{{ end -}},
{{ end -}}
	)
{{ end -}}
{{ end -}}
}

func (p *{{ .ParserName | unexported }}) describe() interface{} {
	{{ if .Fields -}}
	return {{ .ParserName | unexported }}Description {
{{ range $o := .Fields -}}
		{{ .Name }}: p.{{ .Name | unexported }},
{{ end -}}
	}
	{{ else -}}
	return struct{}{}
	{{ end -}}
}

{{ if .Fields -}}
type {{ .ParserName | unexported }}Description struct{
{{ range $o := .Fields -}}
		{{ .Name }} {{ .Type }} ` + "`" + `json:"{{ .Name | snake_case }},omitempty"` + "`" + `
{{ end -}}
}
{{ if .Slice }}
func (d {{ .ParserName | unexported }}Description ) MarshalJSON() ([]byte, error) {
	var comma string
	if d.Comma != 0 {
		comma = string(d.Comma)
	}
	return json.Marshal(struct {
{{ range $o := .Fields -}}
{{ if .Comma -}}
		Comma string ` + "`" + `json:"comma,omitempty"` + "`" + `
{{ else -}}
		{{ .Name }} {{ .Type }} ` + "`" + `json:"{{ .Name | snake_case }},omitempty"` + "`" + `
{{ end -}}
{{ end -}}
	} {
{{ range $o := .Fields -}}
{{ if .Comma -}}
		{{ .Name }}: comma,
{{ else -}}
		{{ .Name }}: d.{{ .Name }},
{{ end -}}
{{ end -}}
	})
}
{{ end -}}
{{ end -}}
`)
)

func main() {
	if err := executeTmpl(uniOptTmpl, "uni_opt.gen.go", types); err != nil {
		log.Fatal(err)
	}

	for _, s := range types {
		if err := executeTmpl(tmpl, snakeCase(s.MethodName)+".gen.go", s); err != nil {
			log.Fatalf("%v: %v", snakeCase(s.MethodName)+".gen.go", err)
		}
	}
}

func executeTmpl(tmpl *template.Template, filename string, data interface{}) error {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return err
	}

	res, err := format.Source(b.Bytes())
	if err != nil {
		f, fErr := os.Create(filename)
		if fErr != nil {
			return err
		}
		defer f.Close()
		_, _ = f.Write(b.Bytes())
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(res))
	return err
}

func tmplWithFuncs(s string) *template.Template {
	return template.Must(template.New("").Funcs(map[string]interface{}{"unexported": unexported, "snake_case": snakeCase}).Parse(s))
}

func unexported(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func snakeCase(s string) string {
	var runes []rune
	for i, r := range []rune(s) {
		if unicode.IsUpper(r) && i != 0 {
			runes = append(runes, '_')
		}
		runes = append(runes, unicode.ToLower(r))
	}
	return string(runes)
}
