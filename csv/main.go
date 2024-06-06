package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Example struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

// 字符串特殊化处理
func handleStr(field string) string {
	field = strings.ReplaceAll(field, `\`, `\\`)
	if strings.Count(field, `'`)%2 == 0 {
		field = strings.ReplaceAll(field, "'", "\\'")
	}
	field = strings.ReplaceAll(field, ",,", ",\\N,")
	if strings.HasSuffix(field, ",") {
		field += "\\N"
	}
	return field
}

func handler(d interface{}) []Example {
	var examples []Example
	data, ok := d.([]string)
	if !ok {
		fmt.Println("输入数据不是切片字符串")
		return examples
	}
	// 正则表达式，匹配逗号但排除在单引号之间的逗号
	reg := regexp.MustCompile(`('[^']+'|[^,]+)`)
	for _, s := range data {
		//fields := strings.Split(s, ",")
		fields := reg.FindAllString(s, -1)
		if len(fields) != reflect.TypeOf(Example{}).NumField() {
			fmt.Println("字段数量不匹配:", s)
			continue
		}
		example := Example{}
		structValue := reflect.ValueOf(&example).Elem()
		for i, fieldValue := range fields {
			field := structValue.Field(i)
			fieldType := field.Kind()
			// 字符串特殊处理
			if fieldType == reflect.String {
				// 将经过处理后的字符串值设置给结构体中的相应字段。
				field.SetString(handleStr(fieldValue))
			}
			// 整数处理
			if fieldType == reflect.Int {
				id, _ := strconv.Atoi(fieldValue)
				field.SetInt(int64(id))
			}
		}
		examples = append(examples, example)
	}
	return examples
}

func main() {
	//data := []string{"1,John,30", "2,Alice,25", "3,Bob,40"}
	data := []string{"1,'Emily',25", "2,'B,enjamin',35", "3,'Ol\\ivia',28"}
	examples := handler(data)
	fmt.Println(examples) // [{1 \'Emily\' 25} {2 \'B,enjamin\' 35} {3 \'Ol\\ivia\' 28}]
}
