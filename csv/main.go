package main

import (
	"fmt"
	"reflect"
)

type Employee struct {
	name   string  `json:"name"`
	age    int     `json:"age"`
	sex    string  `json:"sex"`
	salary float64 `json:"salary"`
}

func handler(i interface{}) {
	tp := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	kd := val.Kind()
	fmt.Printf("tp:%v, val:%v kd:%v\n", tp, val, kd)
	if kd != reflect.Slice {
		return
	}
	// 获取字段数量
	numField := val.NumField()
	fmt.Println("numField:", numField)

	// 遍历结构体
	for i := 0; i < numField; i++ {
		fmt.Printf("field%d 值为%v\n", i, val.Field(i))
	}

}

func main() {
	em := []Employee{
		{name: "张三", age: 22, sex: "男", salary: 11000.5},
		{name: "张三", age: 22, sex: "男", salary: 11000.5},
	}
	handler(em)
}
