package main

import (
	"fmt"
	"strings"
	"unicode"
)

func IsPalindrome(s string) bool {
    // Normalize: lowercase and remove non-letter characters
    var cleaned []rune
    for _, r := range strings.ToLower(s) {
        if unicode.IsLetter(r) || unicode.IsNumber(r) {
            cleaned = append(cleaned, r)
        }
    }

    left, right := 0, len(cleaned)-1
    for left < right {
        if cleaned[left] != cleaned[right] {
            return false
        }
        left++
        right--
    }
    return true
}

func main() {
	one := "hello"
	fmt.Println(IsPalindrome(one))

	two := "wow"
	fmt.Println(IsPalindrome(two))
}