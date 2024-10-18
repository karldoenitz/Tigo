// Package binding 提供了一个非常简单的json以及form的校验功能，支持多种类型的校验。同时也可以自定义函数对某个字段进行校验。主要有两个函数
// 出口`binding.ParseJsonToInstance`和`binding.ValidateInstance`函数。这两个分别用来校验json和form，也可以通过在Param结构体上挂载
// Check函数来自定义校验逻辑。
//
// 定义一个名为UserInfo的Param结构体，包含用户名称、年龄、区域，其中regex这个tag表示该字段必须符合regex指定的正则表达式，如下所示：
//
// Basic Example:
//
//	type UserInfo struct {
//		Name     string `json:"name" form:"name" regex:"^[0-9a-zA-Z_]{1,}$"`
//		Age      int    `json:"age" form:"age"`
//		Location string `json:"location" form:"location"`
//	}
//	ui := UserInfo{}
//	if err := handler.CheckParamBinding(&ui); err != nil {
//		handler.ResponseAsText(err.Error())
//		return
//	}
//
// 如果我们有较为复杂的教研逻辑，比如需要根据区域判断年龄，例如AK地区的年龄必须大于14，MI地区的年龄必须大于16，我们就可以对刚才的UserInfo
// 结构体挂载一个Check校验函数，代码示例如下：
//
// Basic Example:
//
//	func (u *UserInfo) Check() error {
//		if u.Location == "AK" && u.Age <= 14 {
//			return errors.New("AK地区用户年龄必须大于14岁")
//		}
//		if u.Location == "MI" && u.Age <= 16 {
//			return errors.New("MI地区用户年龄必须大于16岁")
//		}
//		return nil
//	}
package binding
