package going

import (
	//go语言库
	//第三访库
	//自己写的工具库
	//项目
	"fmt"
)

//常量统一采用大写，单词之前采用下划线分隔
const (
	CONST_A = "a"
	CONST_B = "b"
)

//变量统一采用驼峰命名格式
var (
	CreatedBy string = "SmartWalle"
)

type IUser interface {
}

type TestUser struct {
	Username string
	Password string
}

func NewTestUser() *TestUser {
	var user = &TestUser{}
	return user
}

func (this *TestUser)Method1() {
	fmt.Println("test")
}

func Func() {
}