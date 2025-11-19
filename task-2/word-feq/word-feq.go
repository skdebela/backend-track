package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) int {
    if s == "" {
        return 0
    }
    return len(strings.Fields(s))
}

func main() {
	text := "This is a test for work count"
	feq := WordCount(text)
	fmt.Println(feq)
}