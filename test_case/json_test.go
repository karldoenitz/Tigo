package test_case

import (
	"fmt"
	"github.com/karldoenitz/Tigo/binding"
	"testing"
)

type Person struct {
	Name   string `json:"name" required:"true"`
	Age    int    `json:"age" required:"true" default:"18"`
	Mobile string `json:"mobile" required:"true" regex:"^1([38][0-9]|14[57]|5[^4])\\d{8}$"`
	Info   string `json:"info" required:"false"`
}

func TestName(t *testing.T) {
	fmt.Println("test case start...")
	test1 := `{"name": "hello", "age": 27, "info": "hello123", "mobile": "15600634756"}`
	person1 := Person{}
	err := binding.ParseJsonToInstance([]byte(test1), &person1)
	if err != nil {
		panic(err.Error())
	}
	test2 := `{"age": 27, "info": "hello123", "mobile": "15600634756"}`
	person2 := Person{}
	er := binding.ParseJsonToInstance([]byte(test2), &person2)
	if er != nil {
		fmt.Println(er.Error())
	} else {
		panic("test failed")
	}
	fmt.Println("test case successfully")
}

func TestAge(t *testing.T) {
	test1 := `{"name": "hello", "age": 27, "info": "hello123", "mobile": "15600634756"}`
	person1 := Person{}
	err := binding.ParseJsonToInstance([]byte(test1), &person1)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(person1.Age)
	test2 := `{"name": "hello", "info": "hello123", "mobile": "15600634756"}`
	person2 := Person{}
	er := binding.ParseJsonToInstance([]byte(test2), &person2)
	if er != nil {
		panic(er.Error())
	}
	fmt.Println(person2.Age)
}

func TestMobile(t *testing.T) {
	test1 := `{"name": "hello", "age": 27, "info": "hello123", "mobile": "15699634756"}`
	person1 := Person{}
	err := binding.ParseJsonToInstance([]byte(test1), &person1)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(person1.Mobile)
	test2 := `{"name": "hello", "age": 27, "info": "hello123"}`
	person2 := Person{}
	er := binding.ParseJsonToInstance([]byte(test2), &person2)
	if er != nil {
		fmt.Println(er.Error())
	} else {
		panic("test case failed")
	}
	fmt.Println(person2.Mobile)
	test3 := `{"name": "hello", "age": 27, "info": "hello123", "mobile": "1560063475601"}`
	person3 := Person{}
	e := binding.ParseJsonToInstance([]byte(test3), &person3)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		panic("test case failed")
	}
	fmt.Println(person3.Mobile)
}

func TestInfo(t *testing.T) {
	test1 := `{"name": "hello", "age": 27, "info": "hello123", "mobile": "15600634756"}`
	person1 := Person{}
	err := binding.ParseJsonToInstance([]byte(test1), &person1)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(person1.Info)
	test2 := `{"name": "hello", "age": 27, "mobile": "15600634756"}`
	person2 := Person{}
	er := binding.ParseJsonToInstance([]byte(test2), &person1)
	if er != nil {
		panic(er.Error())
	}
	fmt.Println(person2.Info)
}
