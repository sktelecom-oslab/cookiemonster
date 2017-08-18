package main

import (
	"fmt"
	"strings"
)

func formatArray(s []string) string {
	return strings.Replace(strings.Trim(fmt.Sprintf("%s", s), "[]"), " ", "\n", -1)
}
