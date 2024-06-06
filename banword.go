package banword

import (
	goahocorasick "github.com/anknown/ahocorasick"
	"github.com/i31pc17/zerowidth"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Detection struct {
	banWords   []string
	allowWords []string
}

type Detected struct {
	OriWord   string
	Word      string
	AllowWord string
	StartPos  int
	EndPos    int
	Length    int
	Allowed   bool
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
	checkText := strings.ToLower(text)
	textRune := []rune(text)
	var removeTexts []removeText
	if len(pattern) > 0 {
		matchText, matchRemoveTexts, err := matchText(pattern, checkText)
		if err == nil {
			checkText = matchText
			removeTexts = matchRemoveTexts
		}
	}

	var detectedList []Detected
	checkTextRune := []rune(checkText)

	m := new(goahocorasick.Machine)
	if err := m.Build(detector.getBanWords()); err != nil {
		return "", nil, err
	}
	bans := m.MultiPatternSearch(checkTextRune, false)

	for _, t := range bans {
		item := Detected{
			Word:     string(t.Word),
			StartPos: t.Pos,
			EndPos:   t.Pos + len(t.Word),
			Length:   len(t.Word),
			Allowed:  false,
		}
		item.OriWord = item.Word
		detectedList = append(detectedList, item)
	}

	if len(detector.allowWords) > 0 && len(detectedList) > 0 {
		if err := m.Build(detector.getAllowWords()); err != nil {
			return "", nil, err
		}
		allows := m.MultiPatternSearch(checkTextRune, false)

		for _, t := range allows {
			startPos := t.Pos
			endPos := t.Pos + len(t.Word)

			for i, d := range detectedList {
				if d.StartPos >= startPos && d.EndPos <= endPos {
					detectedList[i].Allowed = true
					detectedList[i].AllowWord = string(t.Word)
				}
			}
		}
	}

	// 정규식으로 제거된 단어 보정 처리
	if removeTexts != nil && len(removeTexts) > 0 && len(detectedList) > 0 {
		for i, detect := range detectedList {
			for _, re := range removeTexts {
				if re.pos >= detect.EndPos {
					break
				}
				if re.pos <= detect.StartPos {
					detect.StartPos += re.length
				}
				detect.EndPos += re.length
			}
			detect.Length = detect.EndPos - detect.StartPos
			detect.OriWord = string(textRune[detect.StartPos:detect.EndPos])
			detectedList[i] = detect
		}
	}

	if len(detectedList) > 0 {
		for i, d := range detectedList {
			// 공백을 제외하고 단어에 정규식에 제외된 문자가 너무 많이 포함되면 금칙어 제외
			if !d.Allowed && removeTexts != nil {
				wordLen := utf8.RuneCountInString(d.Word)
				// 공백 제거, 안보이는 문자도 제거 해서 검사
				oriWord := strings.ReplaceAll(d.OriWord, " ", "")
				zw := zerowidth.NewZeroWidth()
				zwStr, err := zw.Remove(oriWord)
				if err == nil {
					oriWord = zwStr
				}
				oriWordLen := utf8.RuneCountInString(oriWord)
				if wordLen <= 5 {
					wordLen += 5
				} else {
					wordLen *= 2
				}
				if oriWordLen > wordLen {
					detectedList[i].Allowed = true
				}
			}

			if !d.Allowed {
				for j := d.StartPos; j < d.EndPos; j++ {
					textRune[j] = replaceChar
				}
			}
		}
	}

	return string(textRune), detectedList, nil
}

func (detector *Detection) getBanWords() [][]rune {
	var banWords [][]rune
	for _, word := range detector.banWords {
		banWords = append(banWords, []rune(strings.ToLower(word)))
	}
	return banWords
}

func (detector *Detection) getAllowWords() [][]rune {
	var allowWords [][]rune
	for _, word := range detector.allowWords {
		allowWords = append(allowWords, []rune(strings.ToLower(word)))
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
