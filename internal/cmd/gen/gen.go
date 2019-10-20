package main

import (
	"bytes"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

const (
	global = 1 << iota
	field
	parseParam
)

var (
	placeholder = param{flags: parseParam}

	base         = param{"Base", "int", "", "0", "strconv", field | parseParam}
	bitSize      = param{"BitSize", "int", "", "0", "strconv", field | parseParam}
	layout       = param{"Layout", "string", "time.RFC3339", `""`, "", field | parseParam}
	comma        = param{"Comma", "rune", "", "0", "", field}
	b64Padding   = param{"Padding", "rune", "", "0", "", field | parseParam}
	b64NoPadding = param{"NoPadding", "bool", "", "false", "strconv", field | parseParam}
	b64URLSafe   = param{"URLSafe", "bool", "", "false", "strconv", field | parseParam}

	optional = param{"Optional", "", "", "", "", global}

	types = []specCfg{
		{"Bool", "bool", "strconv.ParseBool", []string{"strconv"}, []param{placeholder}},
		{"Bytes", "[]byte", "parseBytes", nil, []param{placeholder, b64Padding, b64NoPadding, b64URLSafe}},
		{"Duration", "time.Duration", "time.ParseDuration", []string{"time"}, []param{placeholder}},
		{"Float", "float64", "strconv.ParseFloat", []string{"strconv"}, []param{placeholder, bitSize}},
		{"Int", "int64", "strconv.ParseInt", []string{"strconv"}, []param{placeholder, base, bitSize}},
		{"IntSlice", "[]int64", "strconv.ParseInt", []string{"strconv"}, []param{placeholder, base, bitSize}},
		{"String", "string", "", nil, nil},
		{"StringSlice", "[]string", "", nil, nil},
		{"Time", "time.Time", "time.Parse", []string{"time"}, []param{layout, placeholder}},
		{"Uint", "uint64", "strconv.ParseUint", []string{"strconv"}, []param{placeholder, base, bitSize}},
	}
)

type param struct {
	Name, Type, Default, ZeroVal, Import string
	flags                                int
}

func (o param) Global() bool {
	return o.flags&global > 0
}

type specCfg struct {
	MethodName, TypeName string
	ParseFunc            string
	imports              []string
	options              []param
}

func (s specCfg) ParserName() string {
	return s.MethodName + "Parser"
}

func (s specCfg) OptName() string {
	return s.MethodName + "Opt"
}

func (s specCfg) allOptions() []param {
	options := append(make([]param, 0, len(s.options)), s.options...) // copy errytime
	if s.Slice() {
		options = append(options, comma)
	}
	options = append(options,
		param{"Default", s.TypeName, "", "", "", 0},
		optional,
	)
	return options
}

func (s specCfg) Options() []param {
	return s.filterOpts(func(o param) bool {
		return o != placeholder
	}, true)
}

func (s specCfg) LocalOptions() []param {
	return s.filterOpts(func(o param) bool {
		return o != placeholder && !o.Global()
	}, true)
}

func (s specCfg) ParseParams() []param {
	return s.filterOpts(func(o param) bool {
		return o.flags&parseParam > 0
	}, false)
}

func (s specCfg) Fields() []param {
	return s.filterOpts(func(o param) bool {
		return o.flags&field > 0
	}, true)
}

func (s specCfg) filterOpts(f func(o param) bool, doSort bool) (opts []param) {
	for _, o := range s.allOptions() {
		if f(o) {
			opts = append(opts, o)
		}
	}
	if doSort {
		sort.Slice(opts, func(i, j int) bool {
			return opts[i].Name < opts[j].Name
		})
	}
	return
}

func (s specCfg) CustomJSON() bool {
	for _, o := range s.Fields() {
		if o.Type == "rune" {
			return true
		}
	}
	return false
}

func (s specCfg) Slice() bool {
	return strings.Contains(strings.ToLower(s.MethodName), "slice")
}

func (s specCfg) Imports() []string {
	var (
		imports = s.imports
		seen    = make(map[string]bool)
	)
	push := func(s string) {
		if !seen[s] {
			imports = append(imports, s)
			seen[s] = true
		}
	}
	if s.Slice() && s.ParseFunc != "" {
		push("fmt")
	}
	if s.CustomJSON() {
		push("encoding/json")
	}
	for _, o := range s.allOptions() {
		if o.Import != "" {
			push(o.Import)
		}
	}
	push("strings")
	push("errors")
	sort.Strings(imports)
	return imports
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: SPEC_TEMPLATE UNI_OPT_TEMPLATE")
	}
	tmpl, err := fileTmplWithFuncs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	uniOptTmpl, err := fileTmplWithFuncs(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

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

func fileTmplWithFuncs(fName string) (*template.Template, error) {
	return template.New(filepath.Base(fName)).
		Funcs(map[string]interface{}{"unexported": unexported, "snake_case": snakeCase}).
		ParseFiles(fName)
}

func unexported(s string) string {
	if s == "" {
		return s
	}
	s = strings.ReplaceAll(s, "URL", "Url")
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func snakeCase(s string) string {
	var runes []rune
	for i, r := range []rune(strings.ReplaceAll(s, "URL", "Url")) {
		if unicode.IsUpper(r) && i != 0 {
			runes = append(runes, '_')
		}
		runes = append(runes, unicode.ToLower(r))
	}
	return string(runes)
}
