// Code generated by internal/cmd/gen/gen.go DO NOT EDIT.
package envcfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// IntSlice extracts and parses the variable provided according to the options provided.
// Available options:
// - base
// - bit_size
// - comma
// - default
// - optional
func (c *Cfg) IntSlice(docOpts string, opts ...IntSliceOpt) (v []int64) {
	s, err := newIntSliceSpec(docOpts, opts)
	if err != nil {
		if c.panic {
			panic(err)
		}
		c.addError(err)
		return
	}
	c.addDescription(s.describe())
	v, _ = c.evaluate(s).([]int64)
	return
}

// IntSliceOpt modifies IntSlice variable configuration.
type IntSliceOpt interface {
	modify(s *spec)
	modifyIntSliceParser(p *intSliceParser)
}

var IntSlice = struct {
	Base    func(int) IntSliceOpt
	BitSize func(int) IntSliceOpt
	Comma   func(rune) IntSliceOpt
	Default func([]int64) IntSliceOpt
}{
	Base: func(base int) IntSliceOpt {
		return intSliceOptFunc(func(p *intSliceParser) {
			p.setBase(base)
		})
	},
	BitSize: func(bitSize int) IntSliceOpt {
		return intSliceOptFunc(func(p *intSliceParser) {
			p.setBitSize(bitSize)
		})
	},
	Comma: func(comma rune) IntSliceOpt {
		return intSliceOptFunc(func(p *intSliceParser) {
			p.setComma(comma)
		})
	},
	Default: func(def []int64) IntSliceOpt {
		return defaultOpt(def)
	},
}

type intSliceOptFunc func(p *intSliceParser)

func (f intSliceOptFunc) modifyIntSliceParser(p *intSliceParser) {
	f(p)
}

func (intSliceOptFunc) modify(*spec) {}

var _ IntSliceOpt = new(intSliceOptFunc)

func newIntSliceSpec(docOpts string, opts []IntSliceOpt) (*spec, error) {
	parsed, err := parse(docOpts)
	if err != nil {
		return nil, err
	}

	p := new(intSliceParser)
	s := &spec{
		parser:   p,
		typeName: "[]int64",
		name:     parsed.name,
		comment:  parsed.description,
	}

	for _, f := range parsed.fields {
		var (
			opt IntSliceOpt
			err error
		)
		switch strings.ToLower(f[0]) {
		case "base":
			var val int
			val, err = strconv.Atoi(f[1])
			opt = IntSlice.Base(val)
		case "bit_size":
			var val int
			val, err = strconv.Atoi(f[1])
			opt = IntSlice.BitSize(val)
		case "comma":
			value := []rune(f[1])
			if len(value) != 1 {
				err = errors.New("must be only one rune")
				break
			}
			opt = IntSlice.Comma(value[0])
		case "default":
			val := f[1]
			opt = uniOptFunc(func(s *spec) {
				s.flags |= flagDefaultValString | flagDefaultVal
				s.defaultValS = val
			})
		case "optional":
			if f[1] != "" {
				err = errors.New("optional does not take any arguments")
			}
			opt = Optional
		}
		if err != nil {
			return nil, err
		}
		if opt == nil {
			return nil, errors.New("unknown")
		}
		opt.modify(s)
		opt.modifyIntSliceParser(p)
	}

	for _, opt := range opts {
		opt.modify(s)
		opt.modifyIntSliceParser(p)
	}

	if s.flags&flagDefaultValString > 0 {
		if s.defaultVal, err = p.parse(s.defaultValS); err != nil {
			return nil, err
		}
	}

	return s, nil
}

type intSliceParser struct {
	base    int
	bitSize int
	comma   rune
}

func (p *intSliceParser) setBase(base int) {
	p.base = base
}

func (p *intSliceParser) setBitSize(bitSize int) {
	p.bitSize = bitSize
}

func (p *intSliceParser) setComma(comma rune) {
	p.comma = comma
}

func (p *intSliceParser) parse(s string) (interface{}, error) {
	ses, err := parseSlice(s, p.comma)
	if err != nil {
		return nil, err
	}

	vals := make([]int64, len(ses))
	for i, v := range ses {
		el, err := strconv.ParseInt(
			v,
			p.base,
			p.bitSize,
		)
		if err != nil {
			return nil, fmt.Errorf("%v index: %v", i, err)
		}
		vals[i] = el
	}
	return vals, nil

}

func (p *intSliceParser) describe() interface{} {
	return intSliceParserDescription{
		Base:    p.base,
		BitSize: p.bitSize,
		Comma:   p.comma,
	}
}

type intSliceParserDescription struct {
	Base    int  `json:"base,omitempty"`
	BitSize int  `json:"bit_size,omitempty"`
	Comma   rune `json:"comma,omitempty"`
}

func (d intSliceParserDescription) MarshalJSON() ([]byte, error) {
	var comma string
	if d.Comma != 0 {
		comma = string(d.Comma)
	}
	return json.Marshal(struct {
		Base    int    `json:"base,omitempty"`
		BitSize int    `json:"bit_size,omitempty"`
		Comma   string `json:"comma,omitempty"`
	}{
		Base:    d.Base,
		BitSize: d.BitSize,
		Comma:   comma,
	})
}
