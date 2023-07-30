package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const maxWordShow = 10

func Top10(text string) []string {
	splitWords := strings.Fields(text)

	countedWords := countWords(splitWords)

	groupedWords, keys := groupingWords(countedWords)

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	if len(keys) > maxWordShow {
		keys = keys[:maxWordShow]
	}

	return getTop10(groupedWords, keys)
}

func countWords(words []string) map[string]int {
	countedWords := make(map[string]int)

	for i := range words {
		if _, ok := countedWords[words[i]]; !ok {
			countedWords[words[i]] = 0
		}

		countedWords[words[i]]++
	}

	return countedWords
}

func groupingWords(countedWords map[string]int) (map[int][]string, []int) {
	groupedWords := make(map[int][]string)

	for word, cnt := range countedWords {
		groupedWords[cnt] = append(groupedWords[cnt], word)
	}

	keys := make([]int, 0, len(groupedWords))

	for key := range groupedWords {
		keys = append(keys, key)

		sortStrings(groupedWords[key])
	}

	return groupedWords, keys
}

func getTop10(groupedWords map[int][]string, keys []int) []string {
	words := make([]string, 0, len(keys))

	for i := range keys {
		words = append(words, groupedWords[keys[i]]...)
	}

	if len(words) > maxWordShow {
		words = words[:maxWordShow]
	}

	return words
}

func sortStrings(words []string) {
	sort.Slice(words, func(i, j int) bool {
		return words[i] < words[j]
	})
}
