package lineconfig

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// #comment
// some.hierarchy.of.war = content
// some.hierarchy.of.war(int,nic,string,ipv4,ipv6,ipv4net,ipv6net,posint,posfloat,notnull,hex,base58,base64) = content
// some.hierarchy[n].of something
// some.hierarchy.
// .of.war = content # some.hierarchy.of.war
// .for.you = content # some.hierarchy.for.you
// some.data = alpha
// other.data = $some.data # =alpha    // var must exist before
// whatever = "$some.data" # = $some.data

const (
	commentChar    = '#'
	typeCharBegin  = '('
	typeCharEnd    = ')'
	valueChar      = '='
	varChar        = '$'
	quoteChar      = '"'
	escapeChar     = '\\'
	pathSeparator  = '.'
	arrayCharBegin = '['
	arrayCharEnd   = ']'
)

var (
	ErrValidator  = errors.New("wrong validator")
	ErrContent    = errors.New("wrong content")
	ErrPath       = errors.New("wrong path")
	ErrNoContinue = errors.New("cannot continue path")
)

func trim(s string) string {
	return strings.TrimFunc(s, func(r rune) bool { return unicode.IsSpace(r) })
}

func decomment(s string) string {
	if n := strings.IndexRune(s, commentChar); n >= 0 { // Remove comments
		s = s[:n]
	}
	return s
}

func extractName(s string) (name, remainder string) {
	if n := strings.IndexRune(s, typeCharBegin); n >= 0 {
		return trim(s[:n]), s[n:]
	}
	if n := strings.IndexRune(s, valueChar); n >= 0 {
		return trim(s[:n]), s[n:]
	}
	return trim(s), ""
}

func extractValidator(s string) (validator, remainder string, ok bool) {
	if len(s) == 0 {
		return "", "", true
	}
	if s[0] == typeCharBegin {
		if n := strings.IndexRune(s, typeCharEnd); n >= 0 {
			if len(s) < n+1 {
				return "", "", false
			}
			return strings.ToLower(trim(s[1:n])), s[n+1:], true
		}
		return "", "", false
	}
	return "", s, true
}

func parseLine(l string) (path Path, validator string, content Content, err error) {
	var ok bool
	var name string
	var contentS string
	l = decomment(l)
	if len(l) == 0 {
		return nil, "", nil, nil
	}
	name, validator = extractName(l)
	path, err = splitPath(name)
	if err != nil {
		return nil, "", nil, err
	}
	if len(validator) > 0 {
		if path.IsMap() {
			return nil, "", nil, ErrPath
		}
		validator, contentS, ok = extractValidator(validator)
		if !ok {
			return nil, "", nil, ErrValidator
		}
		content, err = extractContent(contentS)
		if err != nil {
			return nil, "", nil, err
		}
		return path, validator, content, nil
	}
	return path, "", nil, nil
}

func appendPath(prev, next Path) Path {
	r := make(Path, 0, len(prev)+len(next))
	for _, e := range prev {
		r = append(r, e)
	}
	for _, e := range next {
		r = append(r, e)
	}
	return r
}

func ParseConfig(r io.Reader) (Tree, error) {
	var i int
	var prevPath Path
	tree := make(Tree)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		i++
		path, validator, content, err := parseLine(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("%d %s: %s", i, path.ToString(), err)
		}
		if len(path) == 0 {
			continue
		}
		if path.IsContinue() {
			if prevPath == nil || !prevPath.IsMap() {
				return nil, fmt.Errorf("%d '%s': %s", i, path.ToString(), ErrNoContinue)
			}
			path = appendPath(prevPath, path)
		} else if path.IsMap() {
			prevPath = path
			continue
		} else {
			prevPath = nil
		}
		if pVal, ok := content.(Path); ok {
			cVal, err := tree.walk(pVal)
			if err != nil {
				return nil, fmt.Errorf("%d '%s': %s", i, content.String(), err)
			}
			content = ContentString(fmt.Sprintf("%v", cVal))
		}
		if validator == "" {
			validator = "string"
		}
		val, err := Validate(validator, content.String())
		if err != nil {
			return nil, fmt.Errorf("%d '%s': %s", i, content.String(), err)
		}
		if _, err = tree.walk(path, val); err != nil {
			return nil, fmt.Errorf("%d '%s': %s", i, path.ToString(), err)
		}
	}
	return tree, scan.Err()
}
