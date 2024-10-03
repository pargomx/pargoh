package dhistorias

import (
	"regexp"
	"strings"
)

func txtQuitarEspaciosYSaltos(txt string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(txt, " "))
}
