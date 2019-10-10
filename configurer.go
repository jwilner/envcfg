package envcfg

import (
	"fmt"
	"os"
)

type Configurable interface {
	Configure(configurer Configurer)
}

func Configure(configurable Configurable, opts ...Option) error {
	c := New(opts...)
	configurable.Configure(c)
	return c.Err()
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

type Cfg struct {
	envFunc func(s string) (string, bool)

	panic    bool
	errors   []error
	errMaker func(err []error) error
}

func (c *Cfg) Has(s string) bool {
	_, ok := c.envFunc(s)
	return ok
}

func (c *Cfg) HasNot(s string) bool {
	return !c.Has(s)
}

func (c *Cfg) Err() error {
	if len(c.errors) > 0 {
		return c.errMaker(c.errors)
	}
	return nil
}

func (c *Cfg) evaluate(s *spec) interface{} {
	v, ok := c.envFunc(s.name)
	if !ok {
		if s.flags&(flagDefaultVal|flagDefaultValString) > 0 {
			return s.defaultVal
		}
		if s.flags&flagOptional == 0 {
			c.addError(fmt.Errorf("%v: variable is required", s.name))
		}
		return nil
	}
	parsed, err := s.parser.parse(v)
	if err != nil {
		c.addError(fmt.Errorf("%v: %v", s.name, err))
	}
	return parsed
}

func (c *Cfg) addError(err error) {
	if c.panic {
		panic(err)
	}
	c.errors = append(c.errors, err)
}

var _ Configurer = new(Cfg)

type spec struct {
	name, typeName string
	flags          int
	defaultVal     interface{}
	defaultValS    string
	parser         parser
	comment        string
}
