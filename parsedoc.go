package envcfg

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

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
