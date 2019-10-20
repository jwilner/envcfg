package envcfg

import (
	"encoding/json"
	"fmt"
	"os"
)

type Description struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Optional bool                   `json:"optional"`
	Default  *DefaultValDescription `json:"default,omitempty"`
	Params   interface{}            `json:"params,omitempty"`
	Comment  string                 `json:"comment,omitempty"`
}

type DefaultValDescription struct {
	Value interface{}
}

func (d DefaultValDescription) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

func New(opts ...Option) *Cfg {
	var defaults = []Option{
		EnvFunc(os.LookupEnv),
		Panic(false),
		ErrMaker(defaultErrMaker),
	}

	c := new(Cfg)
	for _, o := range append(defaults, opts...) {
		o(c)
	}

	return c
}

//go:generate go run internal/cmd/gen/gen.go
type Cfg struct {
	envFunc func(s string) (string, bool)

	descriptions []Description

	panic    bool
	errors   []error
	errMaker func(err []error) error
}

func (c *Cfg) Has(s string) bool {
	c.addDescription(Description{
		Name:     s,
		Type:     "bool",
		Optional: true,
		Default:  &DefaultValDescription{false},
	})
	_, ok := c.envFunc(s)
	return ok
}

func (c *Cfg) HasNot(s string) bool {
	c.addDescription(Description{
		Name:     s,
		Type:     "bool",
		Optional: true,
		Default:  &DefaultValDescription{true},
	})
	_, ok := c.envFunc(s)
	return !ok
}

func (c *Cfg) Err() error {
	if len(c.errors) > 0 {
		return c.errMaker(c.errors)
	}
	return nil
}

func (c *Cfg) Result() ([]Description, error) {
	return c.descriptions, c.Err()
}

func (c *Cfg) addDescription(desc Description) {
	c.descriptions = append(c.descriptions, desc)
}

func (c *Cfg) evaluate(s *spec) interface{} {
	val, err := s.evaluate(c.envFunc)
	if err != nil {
		c.addError(err)
	}
	return val
}

func (c *Cfg) addError(err error) {
	if c.panic {
		panic(err)
	}
	c.errors = append(c.errors, err)
}

type spec struct {
	parser
	name, typeName string
	flags          int
	defaultVal     interface{}
	defaultValS    string
	comment        string
}

func (s *spec) evaluate(envFunc func(string) (string, bool)) (interface{}, error) {
	v, ok := envFunc(s.name)
	if !ok {
		if s.flags&(flagDefaultVal|flagDefaultValString) > 0 {
			return s.defaultVal, nil
		}
		if s.flags&flagOptional == 0 {
			return nil, fmt.Errorf("%v: variable is required", s.name)
		}
		return nil, nil
	}
	return s.parse(v)
}

func (s *spec) describe() Description {
	desc := Description{
		Name:     s.name,
		Type:     s.typeName,
		Optional: s.flags&flagOptional > 0,
		Comment:  s.comment,
		Params:   s.parser.describe(),
	}
	if s.flags&flagDefaultVal > 0 {
		desc.Default = &DefaultValDescription{s.defaultVal}
	}
	return desc
}
