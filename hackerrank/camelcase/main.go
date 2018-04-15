package main

import (
	"fmt"
)

func main() {
	var inputStr string
	fmt.Scanln(&inputStr)
	camelWordCount := 1
	for _, c := range inputStr {
		if c >= 'A' && c <= 'Z' {
			camelWordCount++
		}
	}
	fmt.Println(camelWordCount)
}
