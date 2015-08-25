package main

import (
	//go语言库
	//第三访库
	//自己写的工具库
	//项目
	"fmt"
)

//常量统一采用大写，单词之间采用下划线分隔
const (
	CONST_A = "a"
	CONST_B = "b"
)

//变量统一采用驼峰命名格式
var (
	CreatedBy string = "SmartWalle"
)

////////////////////////////////////////////////////////////////////////////////
//init 函数
func init() {
}

////////////////////////////////////////////////////////////////////////////////
//接口以大写单词“I”开头
type IUser interface {
}

////////////////////////////////////////////////////////////////////////////////
type TestUser struct {
	Username string
	Password string
}

//创建实例的函数
func NewTestUser() *TestUser {
	var user = &TestUser{}
	return user
}

//类方法
func (this *TestUser) Method1() {
	fmt.Println("test")
}

////////////////////////////////////////////////////////////////////////////////
type TestStruct struct {
	Property1 string
	Property2 string
}

func NewTestStruct() *TestStruct {
	var st = &TestStruct{}
	return st
}

func (this *TestStruct) Method1() {
	fmt.Println("test")
}

////////////////////////////////////////////////////////////////////////////////
//其它函数
func Func() {
}
