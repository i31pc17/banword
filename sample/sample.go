package main

import (
	"fmt"
	"github.com/i31pc17/banword"
)

func main() {
	// 금칙어
	banWords := []string{
		"금칙어", "졸라", "money",
	}

	// 허용할 단어
	allowWords := []string{
		"고르곤졸라",
	}

	detector := banword.NewDetector(banWords, allowWords)

	checkText := "금❤️칙❤️❤️어 테스트 중입니다. 여기 고르곤졸라가 졸3 라 졸라 맛있어요. MonEy 필요하나요❤️"
	fmt.Println("텍스트 : ", checkText)

	text, detectedList, err := detector.BanWords(checkText, '*', `[ㄱ-ㅎ가-힣ㅏ-ㅣa-zA-Z]+`)

	if err == nil {
		fmt.Println("필터링 : ", text)

		fmt.Println("필터링 단어")
		for _, d := range detectedList {
			fmt.Printf("단어 : \"%s\", 금칙어 : \"%s\", 허용 여부 : %t, 허용단어 : \"%s\"\n", d.OriWord, d.Word, d.Allowed, d.AllowWord)
		}
	}
}
