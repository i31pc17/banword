# BanWord

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
	"github.com/i31pc17/banword"
	"fmt"
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
```
```
텍스트 :  금❤️칙❤️❤️어 테스트 중입니다. 여기 고르곤졸라가 졸3 라 졸라 맛있어요. MonEy 필요하나요❤️
필터링 :  ********* 테스트 중입니다. 여기 고르곤졸라가 **** ** 맛있어요. ***** 필요하나요❤️
필터링 단어
단어 : "금❤️칙❤️❤️어", 금칙어 : "금칙어", 허용 여부 : false, 허용단어 : ""
단어 : "졸라", 금칙어 : "졸라", 허용 여부 : true, 허용단어 : "고르곤졸라"
단어 : "졸3 라", 금칙어 : "졸라", 허용 여부 : false, 허용단어 : ""
단어 : "졸라", 금칙어 : "졸라", 허용 여부 : false, 허용단어 : ""
단어 : "MonEy", 금칙어 : "money", 허용 여부 : false, 허용단어 : ""
```
