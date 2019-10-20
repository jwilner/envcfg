package envcfg

// Code generated by internal/cmd/gen/gen.go DO NOT EDIT.

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// IntSlice extracts and parses a []int64 variable using the options provided.
//
// The first argument must be a string with beginning with the variable name as expected in the process environment.
// Any other options -- none of which are required -- may either be specified in the remainder of the string or using
// the type-safe IntSliceOpts.
//
// Available options:
// 		- "base" or IntSliceBase
// 		- "bit_size" or IntSliceBitSize
// 		- "comma" or IntSliceComma
// 		- "default" or IntSliceDefault
// 		- "optional" or Optional
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

// IntSliceBase specifies the base to use for a IntSlice variable.
func IntSliceBase(base int) IntSliceOpt {
	return intSliceOptFunc(func(p *intSliceParser) {
		p.base = base
	})
}

// IntSliceBitSize specifies the bit size to use for a IntSlice variable.
func IntSliceBitSize(bitSize int) IntSliceOpt {
	return intSliceOptFunc(func(p *intSliceParser) {
		p.bitSize = bitSize
	})
}

// IntSliceComma specifies the comma to use for a IntSlice variable.
func IntSliceComma(comma rune) IntSliceOpt {
	return intSliceOptFunc(func(p *intSliceParser) {
		p.comma = comma
	})
}

// IntSliceDefault specifies a default value for a IntSlice variable.
func IntSliceDefault(def []int64) IntSliceOpt {
	return defaultOpt(def)
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
			opt = IntSliceBase(val)
		case "bit_size":
			var val int
			val, err = strconv.Atoi(f[1])
			opt = IntSliceBitSize(val)
		case "comma":
			value := []rune(f[1])
			if len(value) != 1 {
				err = errors.New("must be only one rune")
				break
			}
			opt = IntSliceComma(value[0])
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
