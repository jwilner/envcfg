package envcfg

// Code generated by internal/cmd/gen/gen.go DO NOT EDIT.

import (
	"errors"
	"strings"
	"time"
)

// Time extracts and parses a time.Time variable using the options provided.
//
// The first argument must be a string with beginning with the variable name as expected in the process environment.
// Any other options -- none of which are required -- may either be specified in the remainder of the string or using
// the type-safe TimeOpts.
//
// Available options:
// 		- "default" or TimeDefault
// 		- "layout" or TimeLayout
// 		- "optional" or Optional
func (c *Cfg) Time(docOpts string, opts ...TimeOpt) (v time.Time) {
	s, err := newTimeSpec(docOpts, opts)
	if err != nil {
		if c.panic {
			panic(err)
		}
		c.addError(err)
		return
	}
	c.addDescription(s.describe())
	v, _ = c.evaluate(s).(time.Time)
	return
}

// TimeOpt modifies Time variable configuration.
type TimeOpt interface {
	modify(s *spec)
	modifyTimeParser(p *timeParser)
}

// TimeDefault specifies a default value for a Time variable.
func TimeDefault(def time.Time) TimeOpt {
	return defaultOpt(def)
}

// TimeLayout specifies the layout to use for a Time variable.
func TimeLayout(layout string) TimeOpt {
	return timeOptFunc(func(p *timeParser) {
		p.layout = layout
	})
}

type timeOptFunc func(p *timeParser)

func (f timeOptFunc) modifyTimeParser(p *timeParser) {
	f(p)
}

func (timeOptFunc) modify(*spec) {}

var _ TimeOpt = new(timeOptFunc)

func newTimeSpec(docOpts string, opts []TimeOpt) (*spec, error) {
	parsed, err := parse(docOpts)
	if err != nil {
		return nil, err
	}

	p := new(timeParser)
	s := &spec{
		parser:   p,
		typeName: "time.Time",
		name:     parsed.name,
		comment:  parsed.description,
	}

	for _, f := range parsed.fields {
		var (
			opt TimeOpt
			err error
		)
		switch strings.ToLower(f[0]) {
		case "default":
			val := f[1]
			opt = uniOptFunc(func(s *spec) {
				s.flags |= flagDefaultValString | flagDefaultVal
				s.defaultValS = val
			})
		case "layout":

			opt = TimeLayout(f[1])
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
		opt.modifyTimeParser(p)
	}

	for _, opt := range opts {
		opt.modify(s)
		opt.modifyTimeParser(p)
	}

	if s.flags&flagDefaultValString > 0 {
		if s.defaultVal, err = p.parse(s.defaultValS); err != nil {
			return nil, err
		}
	}

	return s, nil
}

type timeParser struct {
	layout string
}

func (p *timeParser) parse(s string) (interface{}, error) {
	layout := p.layout
	if layout == "" {
		layout = time.RFC3339
	}
	return time.Parse(
		layout,
		s,
	)
}

func (p *timeParser) describe() interface{} {
	return timeParserDescription{
		Layout: p.layout,
	}
}

type timeParserDescription struct {
	Layout string `json:"layout,omitempty"`
}
