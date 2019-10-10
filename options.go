package envcfg

import "time"

type Option func(c *Cfg)

func EnvFunc(envFunc func(string) (string, bool)) Option {
	return func(c *Cfg) {
		c.envFunc = envFunc
	}
}

func Panic(b bool) Option {
	return func(g *Cfg) {
		g.panic = b
	}
}

func ErrMaker(f func(errs []error) error) Option {
	return func(g *Cfg) {
		g.errMaker = f
	}
}

func defaultErrMaker(errs []error) error {
	return errs[0]
}

const (
	flagOptional = 1 << iota
	flagDefaultVal
	flagDefaultValString
)

var Optional UniOpt = uniOptFunc(func(s *spec) error {
	s.flags |= flagOptional
	return nil
})

func defaultOpt(defVal interface{}) uniOptFunc {
	return func(s *spec) error {
		s.defaultVal = defVal
		s.flags |= flagDefaultVal
		return nil
	}
}

func Comment(comment string) UniOpt {
	return uniOptFunc(func(s *spec) error {
		s.comment = comment
		return nil
	})
}

var (
	Bool = struct {
		Default func(bool) BoolOpt
	}{
		Default: func(b bool) BoolOpt {
			return defaultOpt(b)
		},
	}

	Float = struct {
		BitSize func(int) FloatOpt
		Default func(float64) FloatOpt
	}{
		BitSize: func(i int) FloatOpt {
			return floatOptFunc(func(p *floatParser) error {
				p.setBitSize(i)
				return nil
			})
		},
		Default: func(f float64) FloatOpt {
			return defaultOpt(f)
		},
	}

	Int = struct {
		BitSize func(int) IntOpt
		Base    func(int) IntOpt
		Default func(int64) IntOpt
	}{
		BitSize: func(i int) IntOpt {
			return intOptFunc(func(p *intParser) error {
				p.setBitSize(i)
				return nil
			})
		},
		Base: func(i int) IntOpt {
			return intOptFunc(func(p *intParser) error {
				p.setBitSize(i)
				return nil
			})
		},
		Default: func(i int64) IntOpt {
			return defaultOpt(i)
		},
	}

	IntSlice = struct {
		BitSize func(int) IntSliceOpt
		Base    func(int) IntSliceOpt
		Comma   func(rune) IntSliceOpt
		Default func([]int64) IntSliceOpt
	}{
		BitSize: func(i int) IntSliceOpt {
			return intSliceOptFunc(func(p *intSliceParser) error {
				p.setBitSize(i)
				return nil
			})
		},
		Base: func(i int) IntSliceOpt {
			return intSliceOptFunc(func(p *intSliceParser) error {
				p.setBase(i)
				return nil
			})
		},
		Comma: func(r rune) IntSliceOpt {
			return intSliceOptFunc(func(p *intSliceParser) error {
				p.setComma(r)
				return nil
			})
		},
		Default: func(i []int64) IntSliceOpt {
			return defaultOpt(i)
		},
	}

	String = struct {
		Default func(string) StringOpt
	}{
		Default: func(s string) StringOpt {
			return defaultOpt(s)
		},
	}

	StringSlice = struct {
		Comma   func(rune) StringSliceOpt
		Default func([]string) StringSliceOpt
	}{
		Comma: func(r rune) StringSliceOpt {
			return stringSliceOptFunc(func(p *stringSliceParser) error {
				p.setComma(r)
				return nil
			})
		},
		Default: func(s []string) StringSliceOpt {
			return defaultOpt(s)
		},
	}

	Time = struct {
		Layout  func(string) TimeOpt
		Default func(time.Time) TimeOpt
	}{
		Layout: func(s string) TimeOpt {
			return timeOptFunc(func(p *timeParser) error {
				p.setLayout(s)
				return nil
			})
		},
		Default: func(t time.Time) TimeOpt {
			return defaultOpt(t)
		},
	}
	Uint = struct {
		BitSize func(int) UintOpt
		Base    func(int) UintOpt
		Default func(uint64) UintOpt
	}{
		BitSize: func(i int) UintOpt {
			return uintOptFunc(func(p *uintParser) error {
				p.setBitSize(i)
				return nil
			})
		},
		Base: func(i int) UintOpt {
			return uintOptFunc(func(p *uintParser) error {
				p.setBitSize(i)
				return nil
			})
		},
		Default: func(i uint64) UintOpt {
			return defaultOpt(i)
		},
	}
)
