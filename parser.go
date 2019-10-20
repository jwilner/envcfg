package envcfg

import (
	"encoding/base64"
	"encoding/csv"
	"errors"
	"net"
	"strings"
)

type parser interface {
	parse(s string) (interface{}, error)
	describe() interface{}
}

func parseSlice(v string, comma rune) ([]string, error) {
	r := csv.NewReader(strings.NewReader(v))
	if comma != 0 {
		r.Comma = comma
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

func parseBytes(s string, padding rune, noPadding, urlSafe bool) ([]byte, error) {
	enc := base64.StdEncoding
	if urlSafe {
		enc = base64.URLEncoding
	}
	if padding != 0 {
		enc = enc.WithPadding(padding)
	}
	if noPadding {
		enc = enc.WithPadding(base64.NoPadding)
	}
	return enc.DecodeString(s)
}

func parseIP(s string) (net.IP, error) {
	parsed := net.ParseIP(s)
	if parsed == nil {
		return nil, errors.New("invalid IP")
	}
	return parsed, nil
}