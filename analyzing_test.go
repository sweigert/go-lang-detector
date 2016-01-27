package langdet

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestCreateProfile(t *testing.T) {
	sampleText := "Begrüßung"
	sampleResult := map[string]int{
		"ß":  1,
		"u":  1,
		"n":  1,
		"Be": 1,
		"rü": 1,
		"ng": 1,
		"e":  1,
		"_B": 1,
		"ßu": 1,
		"r":  1,
		"ü":  1,
		"gr": 1,
		"un": 1,
		"g_": 1,
		"B":  1,
		"eg": 1,
		"üß": 1,
		"g":  2,
	}

	result := CreateOccurenceMap(sampleText, 1)
	ensure.DeepEqual(t, sampleResult, result)
}

func TestCreateProfile2(t *testing.T) {
	sampleText := "Begrüßung"
	sampleResult := map[string]int{
		"g":   2,
		"ü":   1,
		"üß":  1,
		"ßu":  1,
		"egr": 1,
		"ßun": 1,
		"g__": 1,
		"rü":  1,
		"grü": 1,
		"üßu": 1,
		"ng_": 1,
		"B":   1,
		"r":   1,
		"un":  1,
		"_Be": 1,
		"rüß": 1,
		"e":   1,
		"u":   1,
		"_B":  1,
		"eg":  1,
		"Beg": 1,
		"n":   1,
		"__B": 1,
		"ung": 1,
		"ß":   1,
		"Be":  1,
		"gr":  1,
		"ng":  1,
		"g_":  1,
	}

	result := CreateOccurenceMap(sampleText, 2)
	ensure.DeepEqual(t, sampleResult, result)
}

func TestRanking(t *testing.T) {
	sampleText := "AABBCC"
	result := CreateOccurenceMap(sampleText, 5)
	ensure.NotNil(t, result)
	ranking := CreateRankLookupMap(result)
	ensure.NotNil(t, ranking)

	ensure.True(t, ranking["A"] >= 0 && ranking["A"] <= 4)
	ensure.True(t, ranking["B"] >= 0 && ranking["B"] <= 4)
	ensure.True(t, ranking["C"] >= 0 && ranking["C"] <= 4)
}
