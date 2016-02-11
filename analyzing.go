package langdet

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// maxSampleSize represents the maximum number of tokens per sample, low number can
// cause bad accuracy, but better performance.
// -1 for no maximum
var maxSampleSize = 10000
var janitor *regexp.Regexp

var paddings []string

func init() {
	janitor = regexp.MustCompile("P{L}")

	paddings = make([]string, nDepth+2, nDepth+2)
	for i := 0; i <= nDepth+1; i++ {
		paddings[i] = createPadding(i)
	}

}

// Analyze creates the language profile from a given Text and returns it in a Language struct.
func Analyze(text, name string) Language {
	theMap := CreateOccurenceMap(text, nDepth)
	ranked := CreateRankLookupMap(theMap)
	return Language{Name: name, Profile: ranked, OccurrenceMap: theMap}
}

// creates the map [token] rank from a map [token] occurrence
func CreateRankLookupMap(input map[string]int) map[string]int {
	tokens := make([]Token, len(input))
	counter := 0
	for k, v := range input {
		tokens[counter] = Token{Key: k, Occurrence: v}
		counter++
	}
	sort.Sort(ByOccurrence(tokens))
	result := make(map[string]int)
	length := len(tokens)
	locMaxL := maxSampleSize
	if locMaxL < 0 {
		locMaxL = length
	}
	for i := length - 1; i >= 0 && i > length-locMaxL; i-- {
		result[tokens[i].Key] = length - i
	}
	return result
}

// createOccurenceMap creates a map[token]occurrence from a given text and up to a given gram depth
// gramDepth=1 means only 1-letter tokens are created, gramDepth=2 means 1- and 2-letters token are created, etc.
func CreateOccurenceMap(text string, gramDepth int) map[string]int {
	text = cleanText(text)
	tokens := strings.Split(text, " ")
	result := make(map[string]int)
	for _, token := range tokens {
		analyseToken(result, token, gramDepth)
	}
	return result
}

// analyseToken analyses a token to a certain gramDepth and stores the result in resultMap
func analyseToken(resultMap map[string]int, token string, gramDepth int) {
	if len(token) == 0 {
		return
	}
	for i := 1; i <= gramDepth+1; i++ {
		generateNthGrams(resultMap, token, i)
	}
}

// generateNthGrams creates n-gram tokens from the input string and
// adds the mapping from token to its number of occurrences to the resultMap
func generateNthGrams(resultMap map[string]int, text string, n int) {
	padding := paddings[n-1]
	text = padding + text + padding
	upperBound := utf8.RuneCountInString(text) - (n - 1)

	// buffer array for all runes of an ngram
	var ngram = make([]rune, n)
	// current index in the buffer to be written
	var ngramIndex int

	// for each rune, add it to the ngram-buffer
	// and add to the map, if full.
	for p, runeValue := range text {

		if p == upperBound {
			break
		}
		ngram[ngramIndex] = runeValue
		// increment buffer index modulo size
		ngramIndex = (ngramIndex + 1) % n
		// if 0 again, buffer is full.
		if ngramIndex == 0 {
			resultMap[string(ngram)]++
		}
	}
}

// createPadding surrounds text with a padding
func createPadding(length int) string {
	var buffer bytes.Buffer
	padding := "_"
	for i := 0; i < length; i++ {
		buffer.WriteString(padding)
	}
	return buffer.String()
}

// cleanText removes newlines, special characters and numbers from a input text
func cleanText(text string) string {
	return janitor.ReplaceAllString(text, " ")
}
