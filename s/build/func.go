package build

import (
	"regexp"
	"strings"
)

var (
	regexpFirstLetter = regexp.MustCompile(`^([a-zA-Z0-9])`)
	regexpResetSnake  = regexp.MustCompile(`[_]+([a-zA-Z0-9])`)
	regexpTypeName    = regexp.MustCompile(`^([_a-zA-Z0-9\-]+)`)
)

// ToExportName Convert to model field name.
func ToExportName(name string) string {
	return regexpFirstLetter.ReplaceAllStringFunc(
		regexpResetSnake.ReplaceAllStringFunc(name, func(s string) string {
			m := regexpResetSnake.FindStringSubmatch(s)
			return strings.ToUpper(m[1])
		}), func(s string) string {
			m := regexpFirstLetter.FindStringSubmatch(s)
			return strings.ToUpper(m[1])
		},
	)
}
