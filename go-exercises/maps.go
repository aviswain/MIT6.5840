package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	wordCounts := make(map[string]int)
	words := strings.Fields(s)
	
	for _, word := range words {
		wordCounts[word]++
	}
	
	return wordCounts
}

func main() {
	wc.Test(WordCount)
}
