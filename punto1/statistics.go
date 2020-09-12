package main

type Statistics struct {
	WordsCounter map[string]int
	UniqueWordList []string
	Size int
}

func NewStatistics(wordsCounter map[string]int, uniqueWordList []string, size int) *Statistics {
	return &Statistics{
		WordsCounter: wordsCounter,
		UniqueWordList: uniqueWordList,
		Size: size,
	}
}
