package main

import (
	"fmt"
	"strings"
	"regexp"
)

func WordCount(s string) int {
    if s == "" {
        return 0
    }

	reg, _ := regexp.Compile(`[^\w\s]`)
	text := reg.ReplaceAllString(s, "")

    return len(strings.Fields(text))
}

func main() {
	text := "This is a test for work count!"
	feq := WordCount(text)
	fmt.Println(feq)
}