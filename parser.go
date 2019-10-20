package envcfg

import (
	"encoding/csv"
	"errors"
	"strings"
)

type parser interface {
	parse(s string) (interface{}, error)
	describe() interface{}
}

type bitSizer struct {
	bitSize int
}

func (b *bitSizer) setBitSize(value int) {
	b.bitSize = value
}

type baser struct {
	base int
}

func (b *baser) setBase(value int) {
	b.base = value
}

type slicer struct {
	comma rune
}

func (s *slicer) setComma(comma rune) {
	s.comma = comma
}

type layouter struct {
	layout string
}

func (l *layouter) setLayout(s string) {
	l.layout = s
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
