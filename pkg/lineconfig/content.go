package lineconfig

import "strings"

type Content interface {
	String() string
}

type ContentVariable Path

func (content ContentVariable) String() string {
	return Path(content).String()
}

type ContentString string

func (content ContentString) String() string {
	return string(content)
}

func toContentString(s string) (ContentString, error) {
	var escaped bool
	var quoted bool
	r := make([]rune, 0, len(s))
	for n, c := range s {
		if escaped {
			r = append(r, c)
			escaped = false
			continue
		}
		if c == escapeChar {
			escaped = true
			continue
		}
		if c == quoteChar {
			if n == 0 {
				quoted = true
				continue
			}
			if n == len(s)-1 && quoted {
				return ContentString(r), nil
			}
			return "", ErrContent
		}
		r = append(r, c)
	}
	if quoted {
		return "", ErrContent
	}
	return ContentString(r), nil
}

func extractContent(s string) (Content, error) {
	if n := strings.IndexRune(s, valueChar); n >= 0 {
		s = trim(s[n+1:])
		if len(s) > 0 && s[0] == varChar {
			path, err := splitPath(s[1:])
			if err != nil {
				return nil, err
			}
			if path.IsMap() {
				return nil, ErrContent
			}
			return path, nil
		}
		return toContentString(s)
	}
	return nil, nil
}
