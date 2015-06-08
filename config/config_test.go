package config

import (
	"testing"
)

func Test_InitSharedConfig(t *testing.T) {
	var shared = SharedConfig()
	var err = shared.LoadConfig("./test")

	if err != nil {
		t.Error("加载配置文件出错", err)
	}
}

func Test_SharedConfig_CN(t *testing.T) {
	var name = GetString("name", "name")

	if name != "姓名" {
		t.Error("name 字段值应该为 姓名")
	}
}

func Test_SharedConfig_EN(t *testing.T) {
	SetDefaultConfig("en")
	var name = GetString("name", "姓名")
	if name != "name" {
		t.Error("name 字段值应该为 name")
	}
}

func Test_NotExistKey(t *testing.T) {
	var value = GetString("job", "unknown")
	if value != "unknown" {
		t.Error("job 字段值应该为 unknown")
	}

	var _, err = GetValue("job")
	if err == nil {
		t.Error("job 字段不存在")
	}
}

func Test_ConfigExist(t *testing.T) {
	var shared = SharedConfig()
	if shared.ConfigExist("dr") {
		t.Error("配置 dr 不存在")
	}
	if !shared.ConfigExist("cn") {
		t.Error("配置 cn 存在")
	}
}

func Test_KeyExist(t *testing.T) {
	var shared = SharedConfig()
	if shared.KeyExist("dr", "name") {
		t.Error("配置文件 dr 中不存在 name 字段")
	}
	if !shared.KeyExist("cn", "name") {
		t.Error("配置文件 dr 中存在 name 字段")
	}
}