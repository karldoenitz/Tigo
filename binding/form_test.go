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
	formData := ""
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
