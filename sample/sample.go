package main

import (
	"banword"
	"fmt"
)

func main() {
	// 금칙어
	banWords := []string{
		"금칙어", "졸라",
	}

	// 허용할 단어
	allowWords := []string{
		"고르곤졸라",
	}

	detector := banword.NewDetector(banWords, allowWords)

	checkText := "금13칙##어 테33스트 중입@#니다. 여金金金金金기 고르곤졸라가 졸3 라 졸라 맛있어요."
	fmt.Println("텍스트 : ", checkText)

	text, detected, err := detector.BanWords(checkText, '*', `[ㄱ-ㅎ가-힣ㅏ-ㅣa-zA-Z]+`)

	if err == nil {
		fmt.Println("필터링 : ", text)

		fmt.Println("필터링 단어")
		for _, d := range detected {
			fmt.Printf("단어 : \"%s\", 금칙어 : \"%s\", 허용 여부 : %t, 허용단어 : \"%s\"\n", d.OriWord, d.Word, d.Allowed, d.AllowWord)
		}
	}
}
