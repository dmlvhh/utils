package main

import (
	"fmt"
	"strings"
)

func handleStr(inputString string) string {
	inputString = strings.ReplaceAll(inputString, `\`, `\\`)
	if strings.Count(inputString, `'`)%2 == 0 {
		inputString = strings.ReplaceAll(inputString, "'", "\\'")
	}
	inputString = strings.ReplaceAll(inputString, ",,", ",\\N,")
	if strings.HasSuffix(inputString, ",") {
		inputString += "\\N"
	}
	return inputString
}

func main() {
	//inputString := "a,'dsa,das','b,','\\d',"
	inputString := "2,'Benj,amin',35"
	escapedString := handleStr(inputString)
	fmt.Println("原始数据:", inputString)
	fmt.Println("修改后数据:", escapedString)
}
