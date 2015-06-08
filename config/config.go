package config

import (
	"os"
	"path"
	"strings"
	"encoding/json"
	"errors"
)

func init() {
	SharedConfig()
}

var sharedConfig *Config

// NewConfig 创建新的 Config 实例
func NewConfig() (*Config) {
	var config = &Config{}
	return config
}

// SharedConfig 获取共享的 Config 实例
func SharedConfig() (*Config) {
	if sharedConfig == nil {
		sharedConfig = &Config{}
	}
	return sharedConfig
}

func SetDefaultConfig(configName string) bool {
	return sharedConfig.SetDefaultConfig(configName)
}

func GetValue(key string) (interface{}, error) {
	return sharedConfig.GetValue(key)
}

func Get(key string, defaultValue interface{}) interface{} {
	return sharedConfig.Get(key, defaultValue)
}

func GetList(key string, defaultValue []interface{}) []interface{} {
	return sharedConfig.GetList(key, defaultValue)
}

func GetMap(key string, defaultValue map[string]interface{}) map[string]interface{} {
	return sharedConfig.GetMap(key, defaultValue)
}

func GetString(key string, defaultValue string) string {
	return sharedConfig.GetString(key, defaultValue)
}

func GetInt(key string, defaultValue int) int {
	return sharedConfig.GetInt(key, defaultValue)
}

func GetInt32(key string, defaultValue int32) int32 {
	return sharedConfig.GetInt32(key, defaultValue)
}

func GetInt64(key string, defaultValue int64) int64 {
	return sharedConfig.GetInt64(key, defaultValue)
}

func GetFloat(key string, defaultValue float32) float32 {
	return sharedConfig.GetFloat(key, defaultValue)
}

func GetFloat64(key string, defaultValue float64) float64 {
	return sharedConfig.GetFloat64(key, defaultValue)
}

func GetBool(key string, defaultValue bool) bool {
	return sharedConfig.GetBool(key, defaultValue)
}


type Config struct {
	configDir 		string
	configList		[]string
	currentConfig	string
	data 			map[string]map[string]interface{}
}

// LoadConfig 从目录加载数据
// @configDir 配置文件所在目录
func (this *Config) LoadConfig(configDir string) (error) {
	err := this.loadConfigWithDir(configDir)
	return err
}

func (this *Config) loadConfigWithDir(configDir string) error {
	this.configDir = configDir

	dir, err := os.Open(configDir)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	this.data = make(map[string]map[string]interface{})
	this.configList = make([]string, 0, len(names))

	for index, name := range names {
		path := path.Join(configDir, name)

		file, err := os.Stat(path)
		if err != nil {
			continue
		}

		if !file.IsDir() {
			configName := strings.Split(file.Name(), ".")[0]

			if data, err := this.loadConfigWithPath(path); err == nil {
				if index == 0 {
					this.currentConfig = configName
				}
				this.data[configName] = data
				this.configList = append(this.configList, configName)
			} else {
				return err
			}
		}
	}
	return nil
}

func (this *Config) loadConfigWithPath(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&data)

	return data, err
}

// SetDefaultConfig 设定默认的配置文件
// 设定默认的配置文件，如果指定名称的配置文件不存在，则设定为第一个配置。
// @configName 配置名称
func (this *Config) SetDefaultConfig(configName string) bool {
	_, ok := this.data[configName]
	if ok {
		this.currentConfig = configName
	} else if len(this.configList) > 0{
		this.currentConfig = this.configList[0]
		ok = true
	}
	return ok
}

func (this *Config) GetDefaultConfig() string {
	return this.currentConfig
}

func (this *Config) ConfigExist(configName string) bool {
	_, ok := this.data[configName]
	return ok
}

func (this *Config) KeyExist(configName string, key string) bool {
	_, ok := this.data[configName][key]
	return ok
}

// GetValueWithConfig 从指定配置文件读取数据
// 从指定配置文件读取数据，如果指定配置(@configName)或者key(@key)不存在，则返回默认值(@defaultValue)，并且附上对应的错误信息。
// 如果指定配置信息和key都存在，则什么对应的值，错误信息为空。
// @configName 配置名称
// @key key
// @defaultValue 默认值
func (this *Config) GetValueWithConfig(configName string, key string, defaultValue interface{}) (interface{}, error) {
	//data 为空
	if this.data == nil || len(this.data) == 0 {
		return defaultValue, errors.New("读取配置信息失败")
	}

	//如果 configName 为空字符串的时候，则返回默认值
	if len(configName) == 0 {
		return defaultValue, errors.New("配置名称不能为空字符串")
	}

	//如果没有相应的配置
	if !this.ConfigExist(configName) {
		return defaultValue, errors.New("配置[" + configName + "]不存在")
	}

	//如果没有相应的key
	if !this.KeyExist(configName, key) {
		return defaultValue, errors.New("字段[" + key + "]不存在")
	}

	return this.data[configName][key], nil
}

// GetValue 获取指定字段的值
// 如果字段不存在，则返回nil和错误信息；如果字段存在，则返回其值，error为空。
func (this *Config) GetValue(key string) (interface{}, error) {
	return this.GetValueWithConfig(this.currentConfig, key, nil)
}

// Get 获取指定字段的值
// 如果字段不存在，则返回提供的默认值；如果字段存在，则返回其值。
func (this *Config) Get(key string, defaultValue interface{}) interface{} {
	var value, _ = this.GetValueWithConfig(this.currentConfig, key, defaultValue)
	return value
}

func (this *Config) GetList(key string, defaultValue []interface{}) []interface{} {
	var value = this.Get(key, defaultValue)
	if v, ok := value.([]interface{}); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetMap(key string, defaultValue map[string]interface{}) map[string]interface{} {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(map[string]interface{}); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetString(key string, defaultValue string) string {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(string); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetInt(key string, defaultValue int) int {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(int); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetInt32(key string, defaultValue int32) int32 {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(int32); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetInt64(key string, defaultValue int64) int64 {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(int64); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetFloat(key string, defaultValue float32) float32 {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(float32); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetFloat64(key string, defaultValue float64) float64 {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(float64); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetBool(key string, defaultValue bool) bool {
	var value = this.Get(key, defaultValue)
	if v, ok := value.(bool); ok {
		return v
	}
	return defaultValue
}