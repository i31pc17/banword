package banword

import (
	goahocorasick "github.com/anknown/ahocorasick"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Detection struct {
	banWords   []string
	allowWords []string
}

type Detected struct {
	oriWord   string
	word      string
	allowWord string
	startPos  int
	endPos    int
	length    int
	allowed   bool
}

type removeText struct {
	pos    int
	length int
}

func NewDetector(banWords []string, allowWords []string) *Detection {
	detector := &Detection{
		banWords:   banWords,
		allowWords: allowWords,
	}
	return detector
}

func (detector *Detection) BanWords(text string, replaceChar rune, pattern string) (string, []Detected, error) {
	checkText := text
	textRune := []rune(text)
	var removeTexts []removeText
	if len(pattern) > 0 {
		matchText, matchRemoveTexts, err := matchText(pattern, checkText)
		if err == nil {
			checkText = matchText
			removeTexts = matchRemoveTexts
		}
	}

	var detected []Detected
	checkTextRune := []rune(checkText)

	m := new(goahocorasick.Machine)
	if err := m.Build(detector.getBanWords()); err != nil {
		return "", nil, err
	}
	bans := m.MultiPatternSearch(checkTextRune, false)

	for _, t := range bans {
		item := Detected{
			word:     string(t.Word),
			startPos: t.Pos,
			endPos:   t.Pos + len(t.Word),
			length:   len(t.Word),
			allowed:  false,
		}
		item.oriWord = item.word
		detected = append(detected, item)
	}

	if len(detector.allowWords) > 0 && len(detected) > 0 {
		if err := m.Build(detector.getAllowWords()); err != nil {
			return "", nil, err
		}
		allows := m.MultiPatternSearch(checkTextRune, false)

		for _, t := range allows {
			startPos := t.Pos
			endPos := t.Pos + len(t.Word)

			for i, d := range detected {
				if d.startPos >= startPos && d.endPos <= endPos {
					detected[i].allowed = true
					detected[i].allowWord = string(t.Word)
				}
			}
		}
	}

	if removeTexts != nil && len(removeTexts) > 0 && len(detected) > 0 {
		for i, detect := range detected {
			for _, re := range removeTexts {
				if re.pos >= detect.endPos {
					break
				}
				if re.pos <= detect.startPos {
					detect.startPos += re.length
				}
				detect.endPos += re.length
			}
			detect.length = detect.endPos - detect.startPos
			detect.oriWord = string(textRune[detect.startPos:detect.endPos])
			detected[i] = detect
		}
	}

	if len(detected) > 0 {
		for _, d := range detected {
			if !d.allowed {
				for j := d.startPos; j < d.endPos; j++ {
					textRune[j] = replaceChar
				}
			}
		}
	}

	return string(textRune), detected, nil
}

func (detector *Detection) getBanWords() [][]rune {
	var banWords [][]rune
	for _, word := range detector.banWords {
		banWords = append(banWords, []rune(word))
	}
	return banWords
}

func (detector *Detection) getAllowWords() [][]rune {
	var allowWords [][]rune
	for _, word := range detector.allowWords {
		allowWords = append(allowWords, []rune(word))
	}
	return allowWords
}

func matchText(pattern string, input string) (string, []removeText, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", nil, err
	}

	runes := []rune(input)
	matches := re.FindAllString(input, -1)

	var removeTexts []removeText
	lastMatchEnd := 0

	matchIndexes := re.FindAllStringIndex(input, -1)

	for _, matchIndex := range matchIndexes {
		startRune := utf8.RuneCountInString(input[:matchIndex[0]])
		endRune := utf8.RuneCountInString(input[:matchIndex[1]])

		if startRune > lastMatchEnd {
			removeTexts = append(removeTexts, removeText{pos: lastMatchEnd, length: startRune - lastMatchEnd})
		}
		lastMatchEnd = endRune
	}

	if lastMatchEnd < len(runes) {
		removeTexts = append(removeTexts, removeText{pos: lastMatchEnd, length: len(runes) - lastMatchEnd})
	}

	return strings.Join(matches, ""), removeTexts, nil
}
