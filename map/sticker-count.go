package main

import (
	"fmt"
	"math"
)

func calculateMinStickers(word string) int {
	facebookMap := map[string]int{
		"f": 1,
		"a": 1,
		"c": 1,
		"e": 1,
		"b": 1,
		"o": 2,
		"k": 1,
	}
	wordMap := make(map[rune]int)
	for _, char := range word {
		wordMap[char]++
	}
	numberOfSticker := 0

	for key, val := range wordMap {
		numSticker, found := facebookMap[string(key)]
		if !found {
			continue
		}
		needed := math.Ceil(float64(val) / float64(numSticker))
		numberOfSticker = max(numberOfSticker, int(needed))
	}
	return numberOfSticker
}

func main() {
	fmt.Println(calculateMinStickers("coffee kebab")) // Output: 3
	fmt.Println(calculateMinStickers("book"))         // Output: 1
	fmt.Println(calculateMinStickers("ffacebook"))    // Output: 2
}
