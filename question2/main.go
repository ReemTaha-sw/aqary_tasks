package main

import (
	"fmt"
	"sort"
)

func rearrangeString(s string) string {
	charCount := make(map[rune]int)

	// Count the frequency of each character in the string
	for _, char := range s {
		charCount[char]++
	}

	// Create a slice of runes and sort it by frequency in descending order
	sortedChars := make([]rune, 0, len(charCount))
	for char := range charCount {
		sortedChars = append(sortedChars, char)
	}
	sort.Slice(sortedChars, func(i, j int) bool {
		return charCount[sortedChars[i]] > charCount[sortedChars[j]]
	})

	result := make([]rune, len(s))

	// Initialize the result array with characters having the highest frequency
	index := 0
	for _, char := range sortedChars {
		for charCount[char] > 0 && index < len(s) {
			result[index] = char
			index += 2 // Move to every other index
			charCount[char]--
		}
	}

	// If there are any characters left, place them at the remaining positions
	for i := 0; i < len(result); i++ {
		if result[i] == 0 {
			for _, char := range sortedChars {
				if charCount[char] > 0 {
					result[i] = char
					charCount[char]--
					break
				}
			}
		}
	}

	// Check if the rearranged string satisfies the condition
	for i := 1; i < len(result); i++ {
		if result[i] == result[i-1] {
			return ""
		}
	}

	return string(result)
}

func main() {
	fmt.Println(rearrangeString("aab"))   // Output: "aba"
	fmt.Println(rearrangeString("aaab"))  // Output: ""
	fmt.Println(rearrangeString("abc"))   // Output: "abc"
	fmt.Println(rearrangeString(""))        // Output: ""
	fmt.Println(rearrangeString("a"))       // Output: "a"
	fmt.Println(rearrangeString("aa"))      // Output: ""
	fmt.Println(rearrangeString("aaa"))     // Output: ""
	fmt.Println(rearrangeString("aaaaa"))   // Output: ""
}