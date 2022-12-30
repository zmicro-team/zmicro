package main

import (
	"strings"
	"unicode"
)

var (
	// commonInitialisms is a set of common initialisms.
	// source: https://github.com/golang/lint/blob/master/lint.go
	commonInitialisms = map[string]struct{}{
		"ACL":   {},
		"API":   {},
		"ASCII": {},
		"CPU":   {},
		"CSS":   {},
		"DNS":   {},
		"EOF":   {},
		"GUID":  {},
		"HTML":  {},
		"HTTP":  {},
		"HTTPS": {},
		"ID":    {},
		"IP":    {},
		"JSON":  {},
		"LHS":   {},
		"QPS":   {},
		"RAM":   {},
		"RHS":   {},
		"RPC":   {},
		"SLA":   {},
		"SMTP":  {},
		"SQL":   {},
		"SSH":   {},
		"TCP":   {},
		"TLS":   {},
		"TTL":   {},
		"UDP":   {},
		"UI":    {},
		"UID":   {},
		"UUID":  {},
		"URI":   {},
		"URL":   {},
		"UTF8":  {},
		"VM":    {},
		"XML":   {},
		"XMPP":  {},
		"XSRF":  {},
		"XSS":   {},
	}
	defaultReplacer *strings.Replacer
)

func init() {
	initialismForReplacer := make([]string, 0, len(commonInitialisms)*2)
	for s := range commonInitialisms {
		initialismForReplacer = append(initialismForReplacer, s, strings.Title(strings.ToLower(s))) // nolint
	}

	defaultReplacer = strings.NewReplacer(initialismForReplacer...)
}

// Recombine 转换驼峰字符串为用delimiter分隔的字符串, 特殊字符由DefaultInitialisms决定取代
// example: delimiter = '_'
// 空字符 -> 空字符
// HelloWorld -> hello_world
// Hello_World -> hello_world
// HiHello_World -> hi_hello_world
// IDCom -> id_com
// IDcom -> idcom
// nameIDCom -> name_id_com
// nameIDcom -> name_idcom
func Recombine(str string, delimiter byte, enableLint bool) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return ""
	}
	if enableLint {
		str = defaultReplacer.Replace(str)
	}

	var isLastCaseUpper bool
	var isCurrCaseUpper bool
	var isNextCaseUpper bool
	var isNextNumberUpper bool
	var buf = strings.Builder{}

	for i, v := range str[:len(str)-1] {
		isNextCaseUpper = str[i+1] >= 'A' && str[i+1] <= 'Z'
		isNextNumberUpper = str[i+1] >= '0' && str[i+1] <= '9'

		if i > 0 {
			if isCurrCaseUpper {
				if isLastCaseUpper && (isNextCaseUpper || isNextNumberUpper) {
					buf.WriteRune(v)
				} else {
					if str[i-1] != delimiter && str[i+1] != delimiter {
						buf.WriteRune(rune(delimiter))
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(str)-2 && (isNextCaseUpper && !isNextNumberUpper) {
					buf.WriteRune(rune(delimiter))
				}
			}
		} else {
			isCurrCaseUpper = true
			buf.WriteRune(v)
		}
		isLastCaseUpper = isCurrCaseUpper
		isCurrCaseUpper = isNextCaseUpper
	}

	buf.WriteByte(str[len(str)-1])

	return strings.ToLower(buf.String())
}

// SnakeCase 转换驼峰字符串为用'_'分隔的字符串,特殊字符由DefaultInitialisms决定取代
// example2: delimiter = '_' initialisms = DefaultInitialisms
// IDCom -> id_com
// IDcom -> idcom
// nameIDCom -> name_id_com
// nameIDcom -> name_idcom
func SnakeCase(str string, enableLint bool) string {
	return Recombine(str, '_', enableLint)
}

// Kebab 转换驼峰字符串为用'-'分隔的字符串,特殊字符由DefaultInitialisms决定取代
// example2: delimiter = '-' initialisms = DefaultInitialisms
// IDCom -> id-com
// IDcom -> idcom
// nameIDCom -> name-id-com
// nameIDcom -> name-idcom
func Kebab(str string, enableLint bool) string {
	return Recombine(str, '-', enableLint)
}

// SmallCamelCase to small camel case string
// id_com -> idCom
// idcom -> idcom
// name_id_com -> nameIDCom
// name_idcom -> nameIdcom
func SmallCamelCase(fieldName string, enableLint bool) string {
	fieldName = CamelCase(fieldName)
	if enableLint {
		for k := range commonInitialisms {
			if strings.HasPrefix(fieldName, k) {
				return strings.Replace(fieldName, k, strings.ToLower(k), 1)
			}
		}
	}
	return LowTitle(fieldName)
}

// isSeparator reports whether the rune could mark a word boundary.
// TODO: update when package unicode captures more of the properties.
// see strings isSeparator
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case r >= '0' && r <= '9':
			return false
		case r >= 'a' && r <= 'z':
			return false
		case r >= 'A' && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}

	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

// LowTitle 首字母小写
// see strings.Title
func LowTitle(s string) string {
	// Use a closure here to remember state.
	// Hackish but effective. Depends on Map scanning in order and calling
	// the closure once per rune.
	prev := ' '
	return strings.Map(func(r rune) rune {
		if isSeparator(prev) {
			prev = r
			return unicode.ToLower(r)
		}
		prev = r
		return r
	}, s)
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
