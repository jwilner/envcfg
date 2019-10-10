package envcfg

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type parser interface {
	parse(s string) (interface{}, error)
	describe() interface{}
}

type bitSizer struct {
	value int
}

func (b *bitSizer) bitSize() int {
	return b.value
}

func (b *bitSizer) setBitSize(value int) {
	b.value = value
}

func (b *bitSizer) describe() bitSizeDescription {
	return bitSizeDescription{b.value}
}

type bitSizeDescription struct {
	BitSize int `json:"bit_size,omitempty"`
}

type baser struct {
	value int
}

func (b *baser) base() int {
	return b.value
}

func (b *baser) setBase(value int) {
	b.value = value
}

func (b *baser) describe() baseDescription {
	return baseDescription{b.value}
}

type boolParser struct {
}

func (p *boolParser) parse(s string) (interface{}, error) {
	return parseBool(s, p)
}

func (p *boolParser) describe() interface{} {
	return struct{}{}
}

var _ parser = new(boolParser)

func parseBool(s string, _ *boolParser) (interface{}, error) {
	return strconv.ParseBool(s)
}

type uintParser struct {
	intParser
}

func (p *uintParser) parse(s string) (interface{}, error) {
	return strconv.ParseUint(s, p.base(), p.bitSize())
}

var _ parser = new(uintParser)

type intParser struct {
	baser
	bitSizer
}

func (p *intParser) parse(s string) (interface{}, error) {
	return strconv.ParseInt(s, p.base(), p.bitSize())
}

func (p *intParser) describe() interface{} {
	return struct {
		bitSizeDescription
		baseDescription
	}{
		p.bitSizer.describe(),
		p.baser.describe(),
	}
}

type baseDescription struct {
	Base int `json:"base,omitempty"`
}

type floatParser struct {
	bitSizer
}

func (p *floatParser) parse(s string) (interface{}, error) {
	return strconv.ParseFloat(s, p.bitSize())
}

func (p *floatParser) describe() interface{} {
	return p.bitSizer.describe()
}

type durationParser struct {
}

func (durationParser) describe() interface{} {
	return struct{}{}
}

func (durationParser) parse(s string) (interface{}, error) {
	return time.ParseDuration(s)
}

type stringParser struct {
}

func (stringParser) parse(s string) (interface{}, error) {
	return s, nil
}

func (stringParser) describe() interface{} {
	return struct{}{}
}

type timeParser struct {
	layout string
}

func (p *timeParser) setLayout(layout string) {
	p.layout = layout
}

func (p *timeParser) parse(s string) (interface{}, error) {
	layout := p.layout
	if layout == "" {
		layout = time.RFC3339
	}
	return time.Parse(layout, s)
}

func (p *timeParser) describe() interface{} {
	return struct {
		Layout string `json:"layout,omitempty"`
	}{
		p.layout,
	}
}

type slicer struct {
	comma rune
}

func (s *slicer) setComma(comma rune) {
	s.comma = comma
}

func (s *slicer) describe() sliceDescription {
	return sliceDescription{Comma: s.comma}
}

func (s *slicer) parseSlice(v string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(v))
	if s.comma != 0 {
		r.Comma = s.comma
	}
	res, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	switch len(res) {
	case 0:
		return nil, nil
	case 1:
		return res[0], nil
	default:
		return nil, errors.New("at most one line is supported")
	}
}

type sliceDescription struct {
	Comma rune `json:"comma,omitempty"`
}

type stringSliceParser struct {
	slicer
}

func (p *stringSliceParser) parse(s string) (interface{}, error) {
	return p.slicer.parseSlice(s)
}

func (p *stringSliceParser) describe() interface{} {
	return p.slicer.describe()
}

type intSliceParser struct {
	slicer
	intParser
}

func (p *intSliceParser) parse(s string) (interface{}, error) {
	ses, err := p.parseSlice(s)
	if err != nil {
		return nil, err
	}

	vals := make([]int64, len(ses))
	for i, v := range ses {
		el, err := strconv.ParseInt(v, p.base(), p.bitSize())
		if err != nil {
			return nil, fmt.Errorf("%v index: %v", i, err)
		}
		vals[i] = el
	}
	return vals, nil
}

func (p *intSliceParser) describe() interface{} {
	return intSliceDescription{
		p.slicer.describe(),
		p.baser.describe(),
		p.bitSizer.describe(),
	}
}

type intSliceDescription struct {
	sliceDescription
	baseDescription
	bitSizeDescription
}

func (i intSliceDescription) MarshalJSON() ([]byte, error) {
	var s = struct {
		baseDescription
		bitSizeDescription
		Comma string `json:"comma,omitempty"`
	}{
		i.baseDescription,
		i.bitSizeDescription,
		"",
	}
	if i.Comma != 0 {
		s.Comma = string(i.Comma)
	}
	return json.Marshal(s)
}
