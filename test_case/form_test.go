package test_case

import (
	"fmt"
	"github.com/karldoenitz/Tigo/binding"
	"testing"
)

type UserForm struct {
	Name string  `form:"name" required:"true" regex:"^[0-9a-zA-Z_]{1,}$"`
	Age  int     `form:"age" required:"true"`
	VIP  bool    `form:"vip" required:"true"`
	Cash float64 `form:"cash" required:"true" default:"1.02"`
}

func TestForm(t *testing.T) {
	fmt.Println("test start...")
	formData := "name=king%E5%BB%96&age=23&vip=true&cash=12.89"
	userForm := UserForm{}
	err := binding.FormBytesToStructure([]byte(formData), &userForm)
	if err != nil {
		panic(err)
	}
	fmt.Println(userForm.Name)
	fmt.Println(userForm.Age)
	fmt.Println(userForm.VIP)
	fmt.Println(userForm.Cash)
	fmt.Println("test success")
}

func TestParseForm(t *testing.T) {
	fmt.Println("test start...")
	formData := "name=king&age=23&vip=true&cash=12.89"
	userForm := UserForm{}
	err := binding.ParseFormToInstance([]byte(formData), &userForm)
	if err != nil {
		panic(err)
	}
	formData = "name=king%E5%BB%96&age=23&vip=true&cash=12.89"
	err = binding.ParseFormToInstance([]byte(formData), &userForm)
	if err == nil {
		panic("parse name error!")
	}
	fmt.Println(err.Error())
	fmt.Println("test success")
}
