package langdet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

// the depth of n-gram tokens that are created. if nDepth=1, only 1-letter tokens are created
const nDepth = 4

// defaultMinimumConfidence is the minimum confidence that a language-match must have to be returned as detected language
var defaultMinimumConfidence float32 = 0.7

var defaultLanguages = []Language{}

var DefaultDetector = Detector{&defaultLanguages, defaultMinimumConfidence}

// InitWithDefault initializes the default languages with a provided file
// containing Marshalled array of Languages
func InitWithDefault(filePath string) {
	analyzedInput, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Could not open file of default languages: %v", err))
	}
	err = json.Unmarshal(analyzedInput, &defaultLanguages)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshall default languages: %v", err))
	}
}

// Detector has an array of detectable Languages and methods to determine the closest Language to a text.
type Detector struct {
	Languages         *[]Language
	MinimumConfidence float32
}

// NewDetector returns a new Detector without any language.
// It can be used to add languages selectively.
func NewDetector() Detector {
	return Detector{&[]Language{}, defaultMinimumConfidence}
}

// Add language analyzes a text and creates a new Language with given name.
// The new language will be detectable afterwards by this Detector instance.
func (d *Detector) AddLanguageFromText(textToAnalyze, languageName string) {
	if d.Languages == nil {
		newSlice := make([]Language, 0, 0)
		d.Languages = &newSlice
	}
	analyzedLanguage := Analyze(textToAnalyze, languageName)
	updatedList := append(*d.Languages, analyzedLanguage)
	*d.Languages = updatedList
}

// Add language adds a language to the list of detectable languages by this Detector instance.
func (d *Detector) AddLanguage(languages ...Language) {
	if d.Languages == nil {
		s := make([]Language, 0, 0)
		d.Languages = &s
	}
	l := *d.Languages
	for i := range languages {
		l = append(l, languages[i])
	}
	*d.Languages = l
}

// GetClosestLanguage returns the name of the language which is closest to the given text if it is confident enough.
// It returns undefined otherwise. Set detector's MinimumConfidence for customization.
func (d *Detector) GetClosestLanguage(text string) string {
	if d.MinimumConfidence <= 0 || d.MinimumConfidence > 1 {
		d.MinimumConfidence = defaultMinimumConfidence
	}
	if len(*d.Languages) == 0 {
		fmt.Println("no languages configured for this detector")
		return "undefined"
	}
	occ := CreateOccurenceMap(text, nDepth)
	lmap := CreateRankLookupMap(occ)
	c := d.closestFromTable(lmap)

	if len(c) == 0 || c[0].Confidence < asPercent(d.MinimumConfidence) {
		return "undefined"
	}
	return c[0].Name
}

// GetLanguages analyzes a text and returns the DetectionResult of all languages of this detector.
func (d *Detector) GetLanguages(text string) []DetectionResult {
	occ := CreateOccurenceMap(text, nDepth)
	lmap := CreateRankLookupMap(occ)
	results := d.closestFromTable(lmap)
	return results
}

// closestFromTable compares a lookupMap map[token]rank with all languages of this Detector and returns
// an array containing all DetectionResults
func (d *Detector) closestFromTable(lookupMap map[string]int) []DetectionResult {
	res := []DetectionResult{}
	inputSize := len(lookupMap)
	if inputSize > 300 {
		inputSize = 300
	}
	for _, language := range *d.Languages {
		lSize := len(language.Profile)
		maxPossibleDistance := lSize * inputSize
		dist := getDistance(lookupMap, language.Profile, lSize)
		relativeDistance := 1 - float64(dist)/float64(maxPossibleDistance)
		confidence := int(relativeDistance * 100)
		res = append(res, DetectionResult{Name: language.Name, Confidence: confidence})
	}

	sort.Sort(ResByConf(res))
	return res
}

// getDistance calculates the out-of-place distance between two Profiles,
// taking into account only items of mapA, that have a value bigger then 300
func getDistance(mapA, mapB map[string]int, maxDist int) int {
	var result int
	negMaxDist := ((-1) * maxDist)
	for key, rankA := range mapA {
		if rankA > 300 {
			continue
		}
		var diff int
		if rankB, ok := mapB[key]; ok {
			diff = rankB - rankA
			if diff > maxDist || diff < negMaxDist {
				diff = maxDist
			} else if diff < 0 {
				diff = diff * (-1)
			}
		} else {
			diff = maxDist
		}
		result += diff
	}
	return result
}

// asPercentage takes a float and returns its value in percent, rounded to 1%
func asPercent(input float32) int {
	return int(input * 100)
}
