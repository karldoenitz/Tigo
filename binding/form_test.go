package binding

import (
	"fmt"
	"testing"
)

type UserForm struct {
	Name string  `form:"name"`
	Age  int     `form:"age"`
	VIP  bool    `form:"vip"`
	Cash float64 `form:"cash"`
}

func TestForm(t *testing.T) {
	formData := "name=king%E5%BB%96&age=23&vip=true&cash=12.89"
	userForm := UserForm{}
	err := FormBytesToStructure([]byte(formData), &userForm)
	if err != nil {
		panic(err)
	}
	fmt.Println(userForm.Name)
	fmt.Println(userForm.Age)
	fmt.Println(userForm.VIP)
	fmt.Println(userForm.Cash)
}
