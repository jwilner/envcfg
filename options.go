package envcfg

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
