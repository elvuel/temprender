package textual

import (
	"fmt"
	"regexp"
)

type Substitute struct {
	Expr   *string `json:"expr,omitempty"`
	Global bool    `json:"global,omitempty"`
	// only works if `Expr` is nil, the true means prepend replacer otherwise prepend
	Prepend bool `json:"prepend,omitempty"`
	// with format specifier eg: "$0\n%s\n"
	FmtSpecifier string `json:"fmt_specifier,omitempty"`
}

func (sub *Substitute) Sub(src, data string) string {
	repl := fmt.Sprintf(sub.FmtSpecifier, data)

	if sub.Expr == nil {
		if sub.Prepend {
			return repl + src
		} else {
			return src + repl
		}
	}

	pattern := regexp.MustCompile(*sub.Expr)

	if sub.Global {
		return pattern.ReplaceAllString(src, repl)
	}

	first := false
	return pattern.ReplaceAllStringFunc(src, func(matched string) string {
		if first {
			return matched
		}
		first = true
		return pattern.ReplaceAllString(matched, repl)
	})
}
