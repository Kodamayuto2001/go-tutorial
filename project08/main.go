package main

import (
	"fmt"
	"strconv"
)

type Element interface{}
type List [] Element

type Person struct {
	name string
	age int
}

//	Stringメソッドを定義します。fmt.Stringerを実装します。
func (p Person) String() string {
	return "(name: " + p.name + " - age: " + strconv.Itoa(p.age) + " years"
}

func main() {
	list := make(List ,3)
	list[0] = 1
	list[1] = "Hello"
	list[2] = Person{"Dennis",70}

	for index, element := range list {
		if value, ok := element.(int); ok {
			fmt.Printf("list[%d] is an int and its value is %d\n",index,value)
		} else if value, ok := element.(string); ok {
			fmt.Printf("list[%d] is a string and its value is %s\n",index,value)
		} else if value, ok := element.(Person); ok {
			fmt.Printf("list[%d] is a Person and its value is %s\n",index,value)
		} else {
			fmt.Printf("list[%d] is of a different type\n",index)
		}
	}
}