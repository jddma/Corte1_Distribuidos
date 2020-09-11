package main

type Statistics struct {
	WordsCounter map[string]int
	UniqueWordList []string
}

func NewStatistics(wordsCounter map[string]int, uniqueWordList []string) *Statistics {
	return &Statistics{
		WordsCounter: wordsCounter,
		UniqueWordList: uniqueWordList,
	}
}
