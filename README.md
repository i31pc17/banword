# Banword

### 소개
* 텍스트에 금칙어를 검사하고 추출합니다.
* 아호코라식 알고리즘을 이용합니다.
  * https://github.com/anknown/ahocorasick
* 금칙어를 제외할 단어를 지정할 수 있습니다.
* 정규식을 이용해서 원하는 언어만 추출해 검사할 수 있습니다.

### 예제
* sample/sample.go 참고
```go
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

	detectedList := banword.NewDetector(banWords, allowWords)

	checkText := "금13칙##어 테33스트 중입@#니다. 여金金金金金기 고르곤졸라가 졸3 라 졸라 맛있어요."
	fmt.Println("텍스트 : ", checkText)

	text, detected, err := detectedList.BanWords(checkText, '*', `[ㄱ-ㅎ가-힣ㅏ-ㅣa-zA-Z]+`)

	if err == nil {
		fmt.Println("필터링 : ", text)

		fmt.Println("필터링 단어")
		for _, d := range detected {
			fmt.Printf("단어 : \"%s\", 금칙어 : \"%s\", 허용 여부 : %t, 허용단어 : \"%s\"\n", d.OriWord, d.Word, d.Allowed, d.AllowWord)
		}
	}
}
```
```
텍스트 :  금13칙##어 테33스트 중입@#니다. 여金金金金金기 고르곤졸라가 졸3 라 졸라 맛있어요.
필터링 :  ******* 테33스트 중입@#니다. 여金金金金金기 고르곤졸라가 **** ** 맛있어요.
필터링 단어
단어 : "금13칙##어", 금칙어 : "금칙어", 허용 여부 : false, 허용단어 : ""
단어 : "졸라", 금칙어 : "졸라", 허용 여부 : true, 허용단어 : "고르곤졸라"
단어 : "졸3 라", 금칙어 : "졸라", 허용 여부 : false, 허용단어 : ""
단어 : "졸라", 금칙어 : "졸라", 허용 여부 : false, 허용단어 : ""
```
