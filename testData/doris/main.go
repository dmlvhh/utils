package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

var sentences []string
var fullText string
var textLen int
var ts int64 = 1711123200000

func initSentences() {
	content, err := os.ReadFile("bly.txt")
	if err != nil {
		log.Fatal(err)
	}
	splits := strings.Split(string(content), "\n")

	for _, split := range splits {
		tmp := strings.TrimSpace(split)
		if len(tmp) < 2200 {
			continue
		}
		tmp = strings.Replace(tmp, "\"", "", -1)
		tmp = strings.Replace(tmp, ",", "，", -1)
		sentences = append(sentences, tmp)
	}

	for _, sentence := range sentences {
		fullText += sentence
	}
	textLen = len(fullText)
}

func main() {
	initSentences()

	for i := 0; i < 1e2; i++ {

		file, err := os.Create(fmt.Sprintf("doris_1e6_%d.csv", i))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		w := csv.NewWriter(file)
		data := getdata()
		for _, d := range data {
			w.Write(d)
		}
		w.Flush()
		fmt.Println(i)
	}

}

func getdata() [][]string {
	ret := [][]string{}
	for i := 0; i < 1e6; i++ {
		ts += 50
		ret = append(ret, []string{time.Unix(ts/1000, 0).Format("2006-01-02 15:04:05"), getText()})
	}
	return ret
}

func getText() string {
	var bt bytes.Buffer

	for len(bt.String()) < 900 {
		s := randText(rand.Intn(50) + 20)
		bt.WriteString(s)
	}
	return bt.String()

}

func randText(k int) string {
	begin := rand.Intn(textLen - 100)
	ret := fullText[begin:min(begin+k, textLen)]
	return cleanString(ret)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func cleanString(s string) string {
	return strings.Trim(s, string(utf8.RuneError))
}
