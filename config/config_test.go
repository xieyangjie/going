package config

import (
	"fmt"
	"testing"
)

func Test_Init(t *testing.T) {
	var shared = SharedConfig()
	sharedConfig.SetConfigFile("./shared_config.json")
	shared.SetValue("width", 100)
	shared.SetValue("height", 200)
	shared.SaveConfig()

	fmt.Println(shared.GetInt("width", 0))
	fmt.Println(shared.GetInt("height", 0))
	fmt.Println(shared.GetInt("unknow", 10))
}
