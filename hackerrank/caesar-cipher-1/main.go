package main

import (
	"bytes"
	"fmt"
)

func main() {
	var n, k int
	var s string
	fmt.Scanln(&n) // Don't need
	fmt.Scanln(&s)
	fmt.Scanf("%d", &k)

	var b bytes.Buffer
	b.Grow(len(s))
	for _, c := range s {
		b.WriteRune(cipherChar(c, k))
	}
	fmt.Println(b.String())
}

func cipherChar(c rune, k int) rune {
	if c >= 'a' && c <= 'z' {
		a := int(c-'a') + k
		return 'a' + rune(a%26)
	} else if c >= 'A' && c <= 'Z' {
		a := int(c-'A') + k
		return 'A' + rune(a%26)
	}

	return c
}
