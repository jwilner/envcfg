package envcfg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func parseDocOpts(docOpts string) ([]UniOpt, error) {
	parsed, err := parse(docOpts)
	if err != nil {
		return nil, err
	}

	opts := []UniOpt{
		uniOptFunc(func(s *spec) error {
			s.name = parsed.name
			s.comment = parsed.description
			return nil
		}),
	}

	for _, p := range parsed.fields {
		f, err := parseOption(p[0], p[1])
		if err != nil {
			return nil, err
		}
		opts = append(opts, f)
	}

	return opts, nil
}

func parseOption(key, value string) (uniOptFunc, error) {
	switch strings.ToLower(key) {
	case "optional":
		return func(s *spec) error {
			s.flags |= flagOptional
			return nil
		}, nil
	case "base":
		base, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid base: %v", err)
		}
		return func(s *spec) error {
			if s, ok := s.parser.(interface{ setBase(base int) }); ok {
				s.setBase(base)
				return nil
			}
			return errors.New("base is not supported")
		}, nil
	case "bit_size":
		bitSize, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid bit size: %v", err)
		}
		return func(s *spec) error {
			if s, ok := s.parser.(interface{ setBitSize(bitSize int) }); ok {
				s.setBitSize(bitSize)
				return nil
			}
			return errors.New("bitSize is not supported")
		}, nil
	case "comma":
		value := []rune(value)
		if len(value) != 1 {
			return nil, errors.New("comma must be a single rune")
		}
		return func(s *spec) error {
			if s, ok := s.parser.(interface{ setComma(comma rune) }); ok {
				s.setComma(value[0])
				return nil
			}
			return errors.New("comma is not supported")
		}, nil
	case "layout":
		return func(s *spec) error {
			if s, ok := s.parser.(interface{ setLayout(layout string) }); ok {
				s.setLayout(value)
				return nil
			}
			return errors.New("layout is not supported")
		}, nil
	case "default":
		return func(s *spec) error {
			s.flags |= flagDefaultValString | flagDefaultVal
			s.defaultValS = value
			return nil
		}, nil
	}
	return nil, errors.New("unknown option")
}

type parsedOpts struct {
	name, description string
	fields            [][2]string
}

var (
	nameRegexp  = regexp.MustCompile(`^\s*(\S+)`)
	fieldRegexp = regexp.MustCompile(`\s*([^=\s]+)(?:=("(?:\\|[^"]|\")+"|\S+))?`)
)

func parse(doc string) (p parsedOpts, _ error) {
	nameMatch := nameRegexp.FindStringSubmatchIndex(doc)
	if nameMatch == nil {
		return p, errors.New("doc must contain at least name")
	}
	// match group 1, left and right idx
	p.name = doc[nameMatch[2]:nameMatch[3]]
	doc = doc[nameMatch[3]:]

	for {
		fieldMatch := fieldRegexp.FindStringSubmatchIndex(doc)
		if fieldMatch == nil {
			break
		}
		var (
			key    = doc[fieldMatch[2]:fieldMatch[3]]
			offset = fieldMatch[3]
			value  = ""
		)

		if key == "---" || key == "|" {
			doc = doc[offset:]
			break
		} else if key == p.name {
			break
		}

		if fieldMatch[4] != -1 {
			offset = fieldMatch[5]

			value = doc[fieldMatch[4]:fieldMatch[5]]
			if value[0] == '"' && value[len(value)-1] == '"' && value != `"` {
				value = value[1 : len(value)-1]
				value = strings.ReplaceAll(value, `\"`, `"`)
				value = strings.ReplaceAll(value, `\\`, `\`)
			}
		}

		p.fields = append(p.fields, [2]string{key, value})
		doc = doc[offset:]
	}
	p.description = strings.TrimFunc(doc, unicode.IsSpace)
	return
}
