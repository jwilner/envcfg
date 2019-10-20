// Code generated by internal/cmd/gen/gen.go DO NOT EDIT.
package envcfg

type UniOpt interface {
	modify(s *spec) error
	modifyBoolParser(p *boolParser) error
	modifyDurationParser(p *durationParser) error
	modifyFloatParser(p *floatParser) error
	modifyIntParser(p *intParser) error
	modifyIntSliceParser(p *intSliceParser) error
	modifyStringParser(p *stringParser) error
	modifyStringSliceParser(p *stringSliceParser) error
	modifyTimeParser(p *timeParser) error
	modifyUintParser(p *uintParser) error
}

type uniOptFunc func(s *spec) error

func (f uniOptFunc) modify(s *spec) error {
	return f(s)
}

func (uniOptFunc) modifyBoolParser(p *boolParser) error {
	return nil
}

func (uniOptFunc) modifyDurationParser(p *durationParser) error {
	return nil
}

func (uniOptFunc) modifyFloatParser(p *floatParser) error {
	return nil
}

func (uniOptFunc) modifyIntParser(p *intParser) error {
	return nil
}

func (uniOptFunc) modifyIntSliceParser(p *intSliceParser) error {
	return nil
}

func (uniOptFunc) modifyStringParser(p *stringParser) error {
	return nil
}

func (uniOptFunc) modifyStringSliceParser(p *stringSliceParser) error {
	return nil
}

func (uniOptFunc) modifyTimeParser(p *timeParser) error {
	return nil
}

func (uniOptFunc) modifyUintParser(p *uintParser) error {
	return nil
}

var _ UniOpt = new(uniOptFunc)